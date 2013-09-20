package main

import (
	"fmt"
	"io"
	"bytes"
	"math/rand"
	"time"
	"os"
	"strings"
	"math"
	"encoding/hex"
	"crypto/sha256"
)

type Block struct {
	Id string 			// unique Id of block
	Parents []string 	// slice of parent chunks
	Data string 		// encoded source data
	Hash string 		// signed content hash
}

// xor takes two strings of equal length and performs a bitwise XOR on them.
// It returns the resulting string, or the empty string if any error occured.
func xor(s1, s2 string) string {
	
	// Ensure strings are same length
	if len(s1) != len(s2) {
		return ""
	}
	
	// Bitwise XOR strings together
	result := make([]byte, len(s1))
	for i := 0; i < len(s1); i++ {
		result[i] = s1[i] ^ s2[i]
	}
	return string(result)
}

// random takes two ints as bounds for generating a pseudorandom number.
// It returns the number generated.
func random(min, max int) int {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Intn((max + 1) - min) + min
}

// subset takes two string slices and checks whether s1 is a subset of s2.
// {1, 2} is a subset of {1, 2, 3}, but {1, 1} is not a subset of {1, 2, 3}.
// It returns true if s1 is a subset of s2, false if it is not.
func subset(s1, s2 []string) bool {
	set := make(map[string]int)
	for _, value := range s2 {
		set[value] += 1
	}

	for _, value := range s1 {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

// prng is a pseudorandom number generator designed to generate a list of
// chunks to combine into a block. It takes a seed string, degree of
// distrobution, and total number of chunks for the file. It returns a list
// of ints designating chunk indexs to combine.
func prng(s string, k int, n int) []int {

	// Bug: Numbers don't change with different seeds
	// fmt.Println(Prng("b", 5, 500))
	// fmt.Println(Prng("3", 5, 500))

	// Convert seed hash to decimal
	blocks := make([]int, k)
	sBytes, _ := hex.DecodeString(s)
	var seed int64
	for _, element := range sBytes {
		seed += int64(element)
	}

	r := rand.New(rand.NewSource(seed))
	
	// for number of combined blocks
	for i := 0; i < k; i++ {
	
		newItem := true

		block := int(math.Trunc(r.Float64() * float64(n + 1)))
		
		// making sure block is not already in array
		for j := 0; j < k; j++ {
			if blocks[j] == block {
				newItem = false
			}
		}

		if newItem && block <= n {
			blocks[i] = block
		} else {
			i--
		}
	}
	return blocks
}

// ValidateBlocks takes a []Block of remotely pulled blocks from the network
// and compares it to a block chain generated list. Returns an empty list if
// all blocks are valid, or a []Block of corrupt blocks otherwise.
func ValidateBlocks(blocks []Block, signature string) []Block {

	for i := 0; i < len(blocks) - 1; i++ {

		// Generating signed data hash
		h := sha256.New()
		io.WriteString(h, blocks[i].Data + signature)
		new_hash := string(h.Sum(nil))

		if new_hash != blocks[i].Hash {
			fmt.Println("Invalid block: " + hex.EncodeToString([]byte(blocks[i].Id)))
		}
	}
	return []Block{}
}

// Fountain takes an io.Reader from a source file and implements a Luby
// Transform fountain code encoding algorithm over the file. It uses degree
// of distribution, chunk size, and a percent of chunks to generate into
// blocks. 
func Fountain(src io.Reader, dist int, size int, perc float32, signature string) []Block {

	// Still requires signature for signing blocks

	// Convert io.Reader to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	s := buf.String()

	// Setting hash seed
	id := s

	// Chunking source string
	chunks := make([]string, len(s) / size)
	for i := 0; i < len(s) / size; i++ {
		chunks[i] = s[i * size:(i * size) + size]
	}

	// Padding last chunk
	remainder := s[(len(s) - 1) - (len(s) % size):]
	remainder = remainder + string([]byte{1})
	for (size - len(remainder)) > 0{
		remainder = remainder + string([]byte{0})
	}
	chunks = append(chunks, remainder)
	
	// Combining chunks into blocks
	blocks := []Block{}
	for i := 0; i < int(float32(len(chunks)) * perc); i++ {

		// Generate Id hash
		h := sha256.New()
		io.WriteString(h, id)
		id = string(h.Sum(nil))
		
		// Randomly selects first block element
		index := random(0, len(chunks) - 1)
		parents := []string{string(index)}
		data := chunks[index]
		
		// Randomly selects subsequent elements
		for j := 0; j < random(0, dist - 1); j++ {
			index := random(0, len(chunks) - 1)
			parents = append(parents, string(index))
			data = xor(data, chunks[index])
		}

		// Hash content and prepend hash to content
		h = sha256.New()
		io.WriteString(h, data + signature)
		shash := string(h.Sum(nil))

		// Append block to resulting []Block
		blocks = append(blocks, Block{id, parents, data, shash})
	}
	return blocks
}

// DeFountain takes a []Block of valid blocks (validated by ValidateBlocks)
// and recompiles the blocks into the source file as an io.Reader. It
// optionally allows for dynamic reconstruction that has a higher success
// rate, but can increase overall runtime.
func DeFountain(blocks []Block, signature string, dynamic bool) io.Reader {

	// Create map of solved chunks
	found := make(map[string]string)

	// Failed attempt counter
	failed := 0

	// Iterates while still unsolved blocks
	for len(blocks) > 0 {

		// Increments count of fails
		failed++

		//fmt.Println(len(blocks))

		// Pop first block from list
		b := blocks[0]
		blocks = blocks[1:len(blocks)]

		// If solved, add to map of solved
		if len(b.Parents) == 1 {
			found[b.Parents[0]] = b.Data
	
		// Otherwise attempt to solve
		} else {

			// Iterate through parents of block
			parents := []string{}
			for i := 0; i < len(b.Parents); i++ {

				// If parent has been solved, XOR in and remove
				if _, exists := found[b.Parents[i]]; exists {
					b.Data = xor(b.Data, found[b.Parents[i]])
					failed = 0
				
				// Otherwise append the parent back to block's parents
				} else {
					parents = append(parents, b.Parents[i])
				}
			}
			b.Parents = parents

			// Complex resolution: check if block is a subset
			if failed > 2 * len(blocks) && dynamic {

				// Iterate back over []Block
				for i := 0; i < len(blocks); i++ {
					
					// If current block is a subset of another
					if subset(b.Parents, blocks[i].Parents) {
						failed = 0

						// XOR contents of blocks
						blocks[i].Data = xor(b.Data, blocks[i].Data)
						
						// Remove parents from the subset from superset
						parents := blocks[i].Parents
						for j := 0; j < len(b.Parents); j++ {
							for k := 0; k < len(parents); k++ {
								if b.Parents[j] == parents[k] {
									copy(parents[k:], parents[k+1:]) 
									parents = parents[:len(parents)-1] 
									break
								}
							}
						}
						blocks[i].Parents = parents
					}
				}
			}

			// Ensures that block still has parents
			if len(b.Parents) > 0 {
				blocks = append(blocks, b)
			}
		}

		// If decoding fails
		// if failed > 3 * len(blocks) {
		// 	return strings.NewReader("")
		// }

	}
	
	// Converts map of solved elements back into io.Reader
	index := 0
	result := ""
	for {
		value, exists := found[string(index)];
		if !exists {
			break
		} else {
			result += value
			index++
		}
	}

	// Remove padding
	for i := len(result) - 1; i > 0; i-- {
		if string(result[i]) == string([]byte{1}) {
			result = result[:i - 1]
			break
		}
	}

	return strings.NewReader(string(result))
}

// bench is a simple benchmarking tool for encoding and decoding test files
// against test parameters.
func bench(filename string, dist int, size int, perc float32) bool {

	// Load file into io.Reader
	fi, err := os.Open(filename)
	if err != nil { panic(err) }
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	// Encodes and decodes file for benching
	signature := ""
	e := Fountain(fi, dist, size, perc, signature)
	ValidateBlocks(e, signature)
	d := DeFountain(e, signature, false)

	// Save io.Reader to file
	fo, err := os.Create("_" + filename)
	if err != nil { panic(err) }
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	io.Copy(fo, d)

	// Check if failed
	// if io.Reader == ""
	// return false

	// Check if hashes match
	// if string(src_hash) != string(res_hash) {}
	// return false
	
	return true
}

func main() {

	bench("sample.jpg", 5, 1024/4, 5)

}

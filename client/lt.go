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
	Id string
	Parents []string
	Data string
}

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

func random(min, max int) int {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Intn((max + 1) - min) + min
}

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

func prng(s string, k int, n int) []int {
	// seed is a hex hash,
	// k is degree of distribution
	// n is total number of blocks

	//convert seed hash to decimal
	blocks := make([]int, k)
	sBytes, _ := hex.DecodeString(s)
	var seed int64
	for _, element := range sBytes {
		seed += int64(element)
	}

	r := rand.New(rand.NewSource(seed))
	
	for i := 0; i < k; i++ { // for number of combined blocks
	
		newItem := true

		block := int(math.Trunc(r.Float64() * float64(n + 1)))
		
		for j := 0; j < k; j++ { // making sure block is not already in array
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

func ValidateBlocks() {
	
}

func Fountain(src io.Reader, dist int, size int, perc float32) []Block {

	// Convert io.Reader to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	s := buf.String()

	// Setting hash seed
	id := s

	// Chunking source string
	// Still requires trunkation fix
	chunks := make([]string, len(s) / size)
	for i := 0; i < len(s) / size; i++ {
		chunks[i] = s[i * size:(i * size) + size]
	}
	
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

		// Append block to resulting []Block
		blocks = append(blocks, Block{id, parents, data})
	}
	return blocks
}

func DeFountain(blocks []Block) io.Reader {

	// Create map of solved chunks
	found := make(map[string]string)

	// Iterates while still unsolved blocks
	for len(blocks) > 0 {

		fmt.Println(len(blocks))

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
				
				// Otherwise append the parent back to block's parents
				} else {
					parents = append(parents, b.Parents[i])
				}
			}
			b.Parents = parents

			// Complex resolution: check if block is a subset
			// Still requires responsive switching for performance
			if false {

				// Iterate back over []Block
				for i := 0; i < len(blocks); i++ {
					
					// If current block is a subset of another
					if subset(b.Parents, blocks[i].Parents) {

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
	return strings.NewReader(string(result))
}

func main() {

	fmt.Println(Prng("b", 5, 500))
	fmt.Println(Prng("3", 5, 500))
	// Bug: Numbers don't change with different seeds

	filename := "sample.jpg"

	// Load file into io.Reader
	fi, err := os.Open(filename)
	if err != nil { panic(err) }
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	// Encodes and decodes file for testing
	e := Fountain(fi, 5, 1024 / 4, 5)
	d := DeFountain(e)

	// Save io.Reader to file
    fo, err := os.Create("_" + filename)
    if err != nil { panic(err) }
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()
    io.Copy(fo, d)
}

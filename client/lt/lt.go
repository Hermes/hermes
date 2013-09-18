package main

import (
    "fmt"
    "io"
    "bytes"
	"math/rand"
	"time"
)

type Block struct {
	Parents []string
	Data string
}

func xor(s1 string, s2 string) string {
	
	// Ensure strings are same length
	if len(s1) != len(s2) {
		return ""
	}
	
	// XOR strings together
	result := make([]byte, len(s1))
	for i := 0; i < len(s1); i++ {
		result[i] = s1[i] ^ s2[i]
	}
	return string(result)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn((max + 1) - min) + min
}

func Fountain(src io.Reader, dist int, size int, perc float32) []Block {

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	s := buf.String()

	// chunks = []
	chunks := make([]string, len(s) / size)

	// for i in range(len(src) / size):
	for i := 0; i < len(s) / size; i++ {
		// 	chunks.append(src[i * size:(i * size) + size])
		chunks[i] = s[i * size:(i * size) + size]
	}
	
	// blocks = []
	blocks := make([]Block, int(float32(len(chunks)) * perc))
	
	// for j in range(int(len(chunks) * perc)):
	for i := 0; i < int(float32(len(chunks)) * perc); i++ {
		
		// 	index = randint(0, len(chunks)-1)
		index := random(0, len(chunks) - 1)

		// 	b = block([index], chunks[index])
		parents := []string{string(index)}
		data := chunks[index]
		
		// 	for i in range(randint(0, dist - 1)):
		for j := 0; j < random(0, dist - 1); j++ {
		
			// 	index = randint(0, len(chunks) - 1)
			index := random(0, len(chunks) - 1)

			// 	b.parents.append(index)
			parents[j + 1] = string(index)

			// 	b.content = xor(b.content, chunks[index])
			data = xor(data, chunks[index])

		}
		// 	blocks.append(b)
		blocks[i] = Block{parents, data}
	}

	return blocks
}

// func DeFountain(blocks []Block) io.Reader {
// 	return io.Reader
// }

func main() {
	string1 := "WhatAmIDoingHere?"
	string2 := "lakjsdf;lasjdfljd"

	fmt.Println(xor(string1, string2))
}
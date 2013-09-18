package main

import (
    "fmt"
    "io"
    "bytes"
	"math/rand"
	"time"
	"os"
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
	rand.Seed(int64(time.Now().Nanosecond()))
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
	blocks := []Block{}
	
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
			parents = append(parents, string(index))

			// 	b.content = xor(b.content, chunks[index])
			data = xor(data, chunks[index])

		}
		// 	blocks.append(b)
		blocks = append(blocks, Block{parents, data})
	}

	return blocks
}

func DeFountain(blocks []Block) string {
	// found = {}
	found := make(map[string]string)

	// while blocks:
	for len(blocks) > 0 {

		// 	b = blocks.pop(0)
		b := blocks[0]
		blocks = blocks[1:len(blocks)]

		// 	if len(b.parents) == 1:
		if len(b.Parents) == 1 {
			//	found[b.parents[0]] = b.content
			found[b.Parents[0]] = b.Data
		
		// 	else:
		} else {

			// 	parents = []
			parents := []string{}
			// 	for parent in b.parents:
			for i := 0; i < len(b.Parents); i++ {
				// 	if found.has_key(parent):
				if _, exists := found[b.Parents[i]]; exists {
					// 	b.content = xor(b.content, found[parent])
					b.Data = xor(b.Data, found[b.Parents[i]])
				// 	else:
				} else {
					// 	parents.append(parent)
					parents = append(parents, b.Parents[i])
				}
			}
			// 	b.parents = parents
			b.Parents = parents

			// 	if True: # enable if stuck
			if true {
				// 	for c in blocks:
				for i := 0; i < len(blocks); i++ {
					// 	if not Counter(b.parents) - Counter(c.parents): # b is a subset of c
						// 	c.content = xor(b.content, c.content)
						// 	parents = c.parents
						//  for i in b.parents:
							//  parents.remove(i)
						//  c.parents = parents
				}
			}

			// 	if len(b.parents) > 0:
			if len(b.Parents) > 0 {
				// 	blocks.append(b)
				blocks = append(blocks, b)
			}
		}	
		fmt.Println(len(blocks))
	}
	

	for i := 0; i < 100; i++ {
        fmt.Println(found[string(i)])
    }
	return ""
}

func main() {

	fi, err := os.Open("sample.txt")
	if err != nil { panic(err) }
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	string1 := "WhatAmIDoingHere?"
	string2 := "lakjsdf;lasjdfljd"

	fmt.Println(xor(string1, string2))

	e := Fountain(fi, 5, 1024, 10.0)
	DeFountain(e)


}
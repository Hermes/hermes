package client

import (
	"fmt"
	"math/rand"
	"math"
	"encoding/hex"
)

func Prng(s string, k int, n int) []int {
	// seed is a hex hash,
	// k is degree of distribution,
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
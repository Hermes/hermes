package main

import (
    "fmt"
)

func main() {
	string1 := "WhatAmIDoingHere?"
	string2 := "thisisjusta˚∆˙˚∆˙"

	n := len(string1)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = string1[i] ^ string2[i]
	}
	fmt.Printf("%x\n", string(b))

	string3 := string(b)
	c := make([]byte, n)
	for i := 0; i < n; i++ {
		fmt.Println(string1[i] ^ string3[i])
	}
	fmt.Printf("%x\n", string(c))

}
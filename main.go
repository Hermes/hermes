// Hermes, a distributed backup system (DBS)

package main

import (
	"hermes/client"
	"os"
	"io"
)

func main() {
	in, _ := os.Open("../README.md")
	defer in.Close()
	io.Copy(os.Stdout, client.Decompress(client.Compress(in)))
}
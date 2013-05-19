// Hermes, a distributed backup system (DBS)

package main

import (
	"hermes/client"
	"hermes/server"
	"os"
	"flag"
	"fmt"
)

const (
	blockSize = 1048576
	tempDir = "temp"
)

type vault struct {
	Key string
}

func generate(file string) {//, pass string) {
	creds := server.NewCredentials()
	fmt.Println("Keep key secret, and safe.")
	fmt.Println("Vault Key: " + creds.String())
}

func load(file string, pass string) {
	// server code
	fmt.Println("Vault Key: " + " loaded")
}

func update() {
	// server code
	fmt.Println("Vault update successful")
}

func pull(file string) {
	//server code
	/*d := client.Join(files)
	d = client.Decrypt(d, "password")
	d = client.Decompress(d)
	file, _ := os.Create("../bench-files/test1.zip")
	io.Copy(file, d)*/
}

func push(file string) {
	in, _ := os.Open(file)
	defer in.Close()
	i := client.Compress(in)
	i = client.Encrypt(i, "password")
	client.Split(i, blockSize, tempDir)
	// server code
}

func lock() {
	fmt.Println("lock")
}

func main() {

	// check if temp dir exists / make it

	flag.Parse()
	flags := flag.Args()
	if client.ValidateFlags(flags) {
		switch flags[0] { 
			case "update": update()
			case "generate": generate(flags[1])//, flags[2])
			case "load": load(flags[1], flags[2])
			case "pull": pull(flags[1])
			case "push": push(flags[1])
			case "lock": lock()
			default: fmt.Println("Error: Invalid Flags")
		}
	} else {
		fmt.Println("Error: Invalid Flags")
	}

}
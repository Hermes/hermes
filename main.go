// Hermes, a distributed backup system (DBS)

package main

import (
	"hermes/client"
	"hermes/server"
	"os"
	"io"
	"flag"
	"fmt"
	"bytes"
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
    f, err := os.Create(file)
    defer f.Close()
    if err != nil {
        fmt.Println(err)
    }
    n, err := io.WriteString(f, creds.String())
    if err != nil {
        fmt.Println(n, err)
    }
	fmt.Println("Keep key secret, and safe.")
	fmt.Println("Vault Key: " + creds.String())
}

func load(file string) {//, pass string) {
    f, _ := os.Open(file)
	defer f.Close()

	fo, err := os.Create("vault.dat")
    defer fo.Close()
    if err != nil {
        fmt.Println(err)
    } else {
    	io.Copy(fo, f)
    }

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	s := buf.String()
	fmt.Println("Vault Key: " + s + " loaded")
}

func (v vault) update() {
	// server code
	fmt.Println(v.Key)
	fmt.Println("Vault update successful")
}

func (v vault) pull(file string) {
	//server code
	/*d := client.Join(files)
	d = client.Decrypt(d, "password")
	d = client.Decompress(d)
	file, _ := os.Create("../bench-files/test1.zip")
	io.Copy(file, d)*/
}

func (v vault) push(file string) {
	in, _ := os.Open(file)
	defer in.Close()
	i := client.Compress(in)
	i = client.Encrypt(i, "password")
	client.Split(i, blockSize, tempDir)
	// server code
}

func lock() {
	err := os.Remove("vault.dat")
	if err != nil {
		fmt.Println("No active vault to lock")
	} else {
		fmt.Println("Vault has been locked")
	}
}

func main() {

	// check if temp dir exists / make it

	var v vault

	flag.Parse()
	flags := flag.Args()
	if client.ValidateFlags(flags) {

		vaultfile, err := os.Open("vault.dat")
		if err != nil && flags[0] != "generate" && flags[0] != "load" {
	        fmt.Println("Failed to load vault")
	        return
		} else if err == nil {
			defer vaultfile.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(vaultfile)
			s := buf.String()
			v.Key = s
		}

		switch flags[0] { 
			case "generate": generate(flags[1])//, flags[2])
			case "load": load(flags[1]) //, flags[2])
			case "lock": lock()
			case "update": v.update()
			case "pull": v.pull(flags[1])
			case "push": v.push(flags[1])
			default: fmt.Println("Error: Invalid Flags")
		}
	} else {
		fmt.Println("Error: Invalid Flags")
	}

}
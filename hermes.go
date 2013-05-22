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

var verbose bool

const (
	hermesVersion	= 0.1
	hermesBuild		= 1001
	blockSize		= 1048576
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
    _, err = io.WriteString(f, creds.String())
    if err != nil {
        fmt.Println(err)
    }
	fmt.Println("Keep key secret, and safe.")
	fmt.Println("Vault Key: " + creds.String())
}

func load(file string) {//, pass string) {
	vprint("Reading vault file")
    f, _ := os.Open(file)
	defer f.Close()
	vprint("Writing to vault.dat")
	fo, err := os.Create("vault.dat")
    defer fo.Close()
    if err != nil {
        fmt.Println(err)
    } else {
    	io.Copy(fo, f)
    }
}

func (v vault) update() {
	// server code
}

func (v vault) pull(file string) {
	// server code
	vprint("Joining blocks")
	d := client.Join(file)
	vprint("Decrypting file")
	d = client.Decrypt(d, v.Key)
	vprint("Decompressing file")
	d = client.Decompress(d)
	vprint("Saving file")
	result, _ := os.Create(file)
	defer result.Close()
	io.Copy(result, d)
}

func (v vault) push(file string) {
	vprint("Reading file")
	in, _ := os.Open(file)
	defer in.Close()
	vprint("Compressing file")
	i := client.Compress(in)
	vprint("Encrypting file")
	i = client.Encrypt(i, v.Key)
	vprint("Spliting file into blocks")
	client.Split(i, blockSize, file)
	// server code
}

func lock() {
	err := os.Remove("vault.dat")
	if err != nil {
		fmt.Println("Error: No active vault to lock")
	} else {
		fmt.Println("Vault has been locked")
	}
}

func vprint (msg string) {
	if verbose {
		fmt.Println(msg)
	}
}

func main() {

	// Vault instantiation
	var v vault

	// Flag variables
	var help bool
	flag.BoolVar(&help, "h", false, "Add -h for help message")
	flag.BoolVar(&help, "help", false, "Add -h for help message")
	flag.BoolVar(&verbose, "v", false, "Add -v for verbose messages")
	flag.Parse()

	// Help handler
	if help {
		fmt.Println("hermes - An open-source distributed unlimited redundant backup solution")
		fmt.Println("<> required / () optional\n")
		fmt.Println("upgrade\t\t\t\tUpgrades hermes client to latest version")
		fmt.Println("generate <file> (password)\tGenerates new vault file")
		fmt.Println("load <file/address> (password)\tLoads vault as active")
		fmt.Println("update\t\t\t\tUpdates the vault.dat manifest")
		fmt.Println("push <file/folder>\t\tPushes file to network")
		fmt.Println("pull <file/folder>\t\tPulls file from network")
		fmt.Println("lock\t\t\t\tLocks vault and deletes vault.dat")
		return
	}

	// Flag argument handling
	flags := flag.Args()
	if client.ValidateFlags(flags) {

		// Loading vault.dat
		vaultfile, err := os.Open("vault.dat")
		if err != nil && flags[0] != "generate" && flags[0] != "load" && flags[0] != "upgrade" {
	        fmt.Println("Failed to load vault")
	        return
		} else if err == nil {
			defer vaultfile.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(vaultfile)
			s := buf.String()
			v.Key = s
		}

		// Arguement cases
		switch flags[0] { 
			case "generate": generate(flags[1])//, flags[2])
			case "load": load(flags[1]) //, flags[2])
			case "lock": lock()
			case "upgrade": client.Upgrade(hermesVersion, hermesBuild)
			case "update": v.update()
			case "pull": v.pull(flags[1])
			case "push": v.push(flags[1])
			default: fmt.Println("Error: Invalid Flags")
		}

	} else {
		fmt.Println("Error: Invalid Flags")
	}

}
// Hermes, a distributed backup system (DBS)

package main

import (
	//"hermes/client"
	//"hermes/server"
	//"os"
	//"io"
	"flag"
	"fmt"
)

func generate() {
	fmt.Println("generate")
}

func loadVault() {
	fmt.Println("loadVault")
}

func update() {
	fmt.Println("update")
}

func pull() {
	fmt.Println("pull")
}

func push() {
	fmt.Println("push")
}

func ls() {
	fmt.Println("ls")
}

func cd() {
	fmt.Println("cd")
}

func lock() {
	fmt.Println("lock")
}

func flagValidate(flags []string) bool {
	allFlags := []string{"generate", "update", "pull", "push", "ls", "cd", "lock"}

	// Validating one flag is passed
	count := 0
	for _, flag := range flags {
		for _, pflag := range allFlags {
			if (flag == pflag) {
            	count++
        	}
        	if count > 1 {
        		return false
        	}
		}
    }
	return true
}

func flagExists(flags []string, flag string) bool {
    for _, key := range flags {
        if (flag == key) {
            return true
        }
    }
    return false
}

func flagParse(flags []string) string {
    return flags[0]
}

func main() {

	// Flag parsing
	flag.Parse()
	flags := flag.Args()
	if flagValidate(flags) {
		flag := flagParse(flags)
		switch flag { 
			case "update": update()
			case "generate": generate()
			case "pull": pull()
			case "push": push()
			case "ls": ls()
			case "cd": cd()
			case "lock": lock()
			default: fmt.Println("Error: Invalid Flags")
		}
	} else {
		fmt.Println("Error: Invalid Flags")
	}


	// Testing code
	/*
	in, _ := os.Open("../bench-files/test.zip")
	defer in.Close()
	i := client.Compress(in)
	i = client.Encrypt(i, "password")
	files := client.Split(i, 1048576, "temp")
	d := client.Join(files)
	d = client.Decrypt(d, "password")
	d = client.Decompress(d)
	file, _ := os.Create("../bench-files/test1.zip")
	io.Copy(file, d)*/
}
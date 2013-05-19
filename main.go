// Hermes, a distributed backup system (DBS)

package main

import (
	//"hermes/client"
	//"hermes/server"
	"os"
	//"io"
	"flag"
	"fmt"
)

func generate(file string, pass string) {
	fmt.Println("generate")
}

func loadVault() {
	fmt.Println("loadVault")
}

func update() {
	fmt.Println("update")
}

func pull(file string) {
	fmt.Println("pull")
}

func push(file string) {
	fmt.Println("push")
}

func lock() {
	fmt.Println("lock")
}

func flagValidate(flags []string) bool {
	allFlags := []string{"generate", "update", "pull", "push", "lock"}

	// Validating number of flags
	if len(flags) == 0 || len(flags) >=  3 {
		return false
	}

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

    // Validating for specific cases
    switch flags[0] { 
		case "update":
			if len(flags) != 1 {
				return false
			}
		case "generate":
			if len(flags) < 2 || len(flags) > 3 {
				return false
			} else if len(flags) == 2 {
				flags = append(flags, "")
			}
		case "pull":
			if len(flags) != 2 {
				return false
			}
		case "push":
			if len(flags) != 2 {
				return false
			} else if _, err := os.Stat(flags[1]); os.IsNotExist(err) {
				return false
			}
		case "lock":
			if len(flags) != 1 {
				return false
			}
		default: return false
	}
	return true
}

func main() {

	flag.Parse()
	flags := flag.Args()
	if flagValidate(flags) {
		switch flags[0] { 
			case "update": update()
			case "generate": generate(flags[1], flags[2])
			case "pull": pull(flags[1])
			case "push": push(flags[1])
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
// Hermes, a distributed backup system (DBS)

package main

import (
	"hermes/client"
	"os"
	//"io"
)

func main() {
	in, _ := os.Open("../test.zip")
	defer in.Close()
	i := client.Compress(in)
	i = client.Encrypt(i, "password")
	client.Split(i, 1048576, "temp")

}
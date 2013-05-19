// Hermes, a distributed backup system (DBS)
// Compression benchmarking

package main

import (
	"hermes/client"
	"os"
	"fmt"
	"strconv"
)

func fileSize(stat os.FileInfo) int {
	var size int
	size = int(stat.Size())
	return size
}

func main() {

	benchdir := "../bench-files"
	blockSize := 1048576

	dir, _ := os.Open(benchdir)
	defer dir.Close()
	files, _ := dir.Readdir(0)
		
	fmt.Println("Block Size: " + strconv.Itoa(blockSize))
	fmt.Println("-------------------")
	
	for _, file := range files {
		in, _ := os.Open(benchdir + "/" + file.Name())
		i := client.Compress(in)
		i = client.Encrypt(i, "password")
		client.Split(i, blockSize, "temp")
		var beforeSize int
		beforeSize = int(file.Size())

		// for file in "temp" get fileSize
		tempdir, _ := os.Open("temp")
		tempfiles, _ := tempdir.Readdir(0)
		afterSize := 0
		for _, temp := range tempfiles {
	    	afterSize += fileSize(temp)
	    	os.Remove("temp/" + temp.Name())
		}
		tempdir.Close()
		fmt.Println(file.Name() + "\t" + strconv.Itoa(beforeSize) + "\t-->\t" + strconv.Itoa(afterSize) + "\t")
		in.Close()
	}

}
// Hermes, a distributed backup system (DBS)
// Chunk size benchmarking

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

	files := []string{"zip", "pdf", "mp3", "doc"}
	sizes := []int{1024, 524288, 1048576}

	for _, size := range sizes {
		
		fmt.Println("Block Size: " + strconv.Itoa(size))
		fmt.Println("---------------------")
		
		for _, file := range files {
			
			in, _ := os.Open("../bench-files/test." + file)
			i := client.Compress(in)
			i = client.Encrypt(i, "password")
			client.Split(i, size, "temp")
			stat, _ := in.Stat()
			var beforeSize int
			beforeSize = int(stat.Size())

			// for file in "temp" get fileSize
			dir, _ := os.Open("temp")
			defer dir.Close()
			files, _ := dir.Readdir(0)
			afterSize := 0
			for _, file := range files {
		    	afterSize += fileSize(file)
		    	os.Remove("temp/" + file.Name())
			}

			fmt.Println(file + "\t" + strconv.Itoa(beforeSize) + "\t-->\t" + strconv.Itoa(afterSize) + "\t" + strconv.Itoa((beforeSize / afterSize)*100) + "%")
			in.Close()
		}
		fmt.Println("")
	}

}
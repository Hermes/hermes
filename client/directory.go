package client

import (
	"log"
	"os"
	"path"
)

func merge(a, b []string) []string {
   result := make([]string, len(a) + len(b))
   copy(result, a)
   copy(result[len(a):], b)
   return result
}

func DirWalk(dirPath string) []string {

	// Opening directory and error checking
	filePaths := make([]string, 0)
	dir, err := os.Open(dirPath)
	dirStat, _ := dir.Stat()
	if os.IsPermission(err) { // Checks for editing permissions
		return filePaths // Change to handle file permission denied error
	} else if !dirStat.IsDir() { // Ensures input is directory
		return []string{dirPath}
	}
	defer dir.Close()

	// Reading contents of directory
	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	// Parsing contents of directory
	for _, file := range files {
		curPath := path.Join(dirPath, file.Name())
		if file.IsDir() { // Recursive directory parsing
			subDir := DirWalk(curPath)
			filePaths = merge(filePaths, subDir)
		} else { // Appending files from directory
			filePaths = append(filePaths, string(curPath))
		}
	}
	return filePaths
}
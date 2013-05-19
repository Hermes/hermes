package client

import (
	"log"
	"os"
)

//Error handler
func handleError(_e error) {
	if _e != nil {
		log.Fatal(_e)
	}
}

//walks the selected folder and returns an array of files as strings
func DirWalk(dirPath string) []string {
	filePaths := make([]string, 0)
	dir, err := os.Open(dirPath)
	//check to see if I have permissions to edit the file
	handleError(err) //change this to handle the file permission denied error
	defer dir.Close()
	fis, err := dir.Readdir(0)
	handleError(err)
	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			//walking through files in the current path and adding them to the array one by one
			for _, newfile := range DirWalk(curPath) {
				filePaths = append([]string(filePaths), newfile)
			}
		} else {
			filePaths = append([]string(filePaths), string(curPath))
		}
	}
	return filePaths
}
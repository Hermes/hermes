package client

import (
	//"fmt"
	"os"
	"log"
	//"strings"
)

//Error handler
func handleError( _e error ) {
      if _e != nil {
        log.Fatal( _e )
      }
    }

//walks the selected folder and returns an array of files as strings
func DirWalk(dirPath string) []string {
	filePaths := make([]string, 0)
	dir, err :=os.Open(dirPath)
	//fmt.Println("%v\n", dir)
	handleError(err)
	defer dir.Close()
	fis, err := dir.Readdir(0)
    handleError(err)
    for _, fi := range fis {
    	curPath := dirPath + "/" + fi.Name()
        if fi.IsDir() {
        	DirWalk(curPath)
        } else {
        	filePaths = append([]string(filePaths), string(curPath))
        }
	}
	return filePaths
}
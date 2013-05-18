package client

import (
	"fmt"
	"os"
	"log"
	"strings"
)

func handleError( _e error ) {
      if _e != nil {
        log.Fatal( _e )
      }
    }

func DirWalk(dirPath string) []string {
	filePaths := make([]string, 0)
	dir, err :=os.Open(dirPath)
	fmt.Printf("%v\n", dir)
	handleError(err)
	defer dir.Close()
	fis, err := dir.Readdir(0)
    handleError(err)
    for _, fi := range fis {
    	curPath := dirPath + "/" + fi.Name()
        if fi.IsDir() {
        	DirWalk(curPath)
        } else {
        	filePaths = append(strings.Join(filePaths, string(curPath)))
        }
	}
	return filePaths
}
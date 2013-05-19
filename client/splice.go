package client

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func randomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	return string(bytes)
}

//Split: Splits a 'file' into blocks of size 'block_size' and outputs them into
//The directory given by 'filedir'
func Split(file io.Reader, block_size int, filename string) string {
	final_files := make([]string, 0)
	buf := make([]byte, block_size)
	//Open the file
	file = bufio.NewReader(file)
	filedir := path.Join(os.TempDir(), filename)
	os.Mkdir(filedir, 0775)
	i := 0
	for {
		i++
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		//Hash the file with SHA256 to determine a unique filename
		h := sha256.New()
		s := string(buf[:n]) + filedir + randomString(8)
		io.WriteString(h, s)
		//Creating the filename by appending filedir with the SHA256 Hash
		filename := path.Join(filedir, hex.EncodeToString(h.Sum(nil))+"."+strconv.Itoa(i))

		//Write the file and add it's relative path to the list
		fo, err := os.Create(filename)
		// fo, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		// make a write buffer
		w := bufio.NewWriter(fo)
		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
		if err = w.Flush(); err != nil {
			panic(err)
		}
		final_files = append(final_files, filename)
	}
	return filedir
}

//Join: Takes a list of 'files' and joins them back together according to it's order
func Join(filedir string) io.Reader {
	dir, _ := os.Open(filedir)
	defer dir.Close()
	files, _ := dir.Readdir(0)
	result := make([]byte, 0)
	for _, file := range files {
		//read the entire file
		buf, err := ioutil.ReadFile(path.Join(filedir, file.Name()))
		if err != nil {
			panic(err)
		}

		//combine the bytes together
		for _, bite := range buf {
			result = append(result, bite)
		}
	}

	//create an io.Reader and return it
	sresult := string(result)
	return strings.NewReader(sresult)
}

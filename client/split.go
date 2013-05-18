package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func write_file(content []byte, filename string) {
	// open output file
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
	if _, err := w.Write(content); err != nil {
		panic(err)
	}
	if err = w.Flush(); err != nil {
		panic(err)
	}
}

func split(file io.Reader, block_size int, filedir string) []string {
	// make a buffer to keep chunks that are read
	files := make([]string, 0)
	buf := make([]byte, block_size)
	file = bufio.NewReader(file)
	for {
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			break
		}
		h := sha256.New()
		s := string(buf[:n])
		io.WriteString(h, s)
		filename := hex.EncodeToString(h.Sum(nil))
		filename = path.Join(filedir, filename)
		write_file(buf[:n], filename)
		files = append(files, filename)
	}
	return files
}

func randomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	return string(bytes)
}

func main() {
	// open input file
	fi, err := os.Open("server/splits.png")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	// files := split(bufio.NewReader(fi), 128)
	for k, v := range split(fi, 128, "tmp") {
		fmt.Printf("key=%v, value=%v\n", k, v)
	}
}

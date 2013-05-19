package client

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"strings"
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

func open_file(filename string) io.Reader {
	// open input file
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	return fi

}
func Split(file io.Reader, block_size int, filedir string) []string {
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

func Join(files []string, block_size int) io.Reader {
	result := make([]byte, 0)
	for _, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		// close fi on exit and check for its returned error
		defer func() {
			if err := fi.Close(); err != nil {
				panic(err)
			}
		}()

		r := bufio.NewReader(fi)
		buf := make([]byte, block_size)
		n, err := r.Read(buf)
		for _, bite := range buf[:n] {
			result = append(result, bite)
		}
	}

	sresult := string(result)
	return strings.NewReader(sresult)
}

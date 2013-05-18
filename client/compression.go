// Requires https://code.google.com/p/lzma/

package client

import (
	"code.google.com/p/lzma"
	"io"
	"os"
)

func compress(file string) {

	// Open file input reader
	in, _ := os.Open(file)
	defer in.Close()

	// Create file output writer
	out, _ := os.Create(file + "lz")
	defer out.Close()

	// Creat LZMA writer
	w := lzma.NewWriterLevel(out, 9) // Highest compression
	defer w.Close()

	// Reading from file as buffer and writing compressed
    buf := make([]byte, 1024)
    for {
        n, err := in.Read(buf)
        if err != nil && err != io.EOF { panic(err) }
        if n == 0 { break }
        w.Write(buf[:n])
    }
}

func decompress(file string) io.Reader {

	// Open file input reader
	in, _ := os.Open(file)
	defer in.Close()

	// Open LZMA reader
	r := lzma.NewReader(in)
	return r
}
// Requires https://code.google.com/p/lzma/

package client

import (
	"code.google.com/p/lzma"
	"io"
	"bytes"
)

func Compress(in io.Reader) io.Reader {

	var out bytes.Buffer

	// Creat LZMA writer
	w := lzma.NewWriterLevel(&out, 9) // Highest compression
	defer w.Close()

	// Reading from file as buffer and writing compressed
    buf := make([]byte, 1024)
    for {
        n, err := in.Read(buf)
        if err != nil && err != io.EOF {
        	panic(err)
        }
        if n == 0 { break }
        w.Write(buf[:n])
    }
	return &out
}

func Decompress(in io.Reader) io.Reader {
	return lzma.NewReader(in)
}
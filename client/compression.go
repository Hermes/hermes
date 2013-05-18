// Requires https://code.google.com/p/lzma/

package main

import (
	//"fmt"
	"code.google.com/p/lzma"
	"io"
	"os"
)

func compress(file string) {
	in, _ := os.Open(file)
	defer in.Close()
	out, _ := os.Create(file + "lz")
	defer out.Close()
	w := lzma.NewWriter(out)

	// make a buffer to keep chunks that are read
    buf := make([]byte, 1024)
    for {
        // read a chunk
        n, err := in.Read(buf)
        if err != nil && err != io.EOF { panic(err) }
        if n == 0 { break }

        w.Write(buf[:n])
        /* write a chunk
        if _, err := out.Write(buf[:n]); err != nil {
            panic(err)
        }*/
    }
	w.Close()
}

func decompress(file string) {
	in, _ := os.Open(file)
	defer in.Close()
	r := lzma.NewReader(in)
	io.Copy(os.Stdout, r)
	r.Close()
}

func main() {
	
	compress("test.mp4")
	//decompress("test.doclz")

	/*

	// If the data is bigger than you'd like to hold into memory, use pipes. Write compressed data to an io.PipeWriter:
	pr, pw := io.Pipe()
	go func() {
	    defer pw.Close()
	    w := lzma.NewWriter(pw)
	    defer w.Close()
	    // the bytes.Buffer would be an io.Reader used to read uncompressed data from
	    io.Copy(w, bytes.NewBuffer([]byte("hello, world\n")))
	}()


	// and read it back:
	defer pr.Close()
	r := lzma.NewReader(pr)
	defer r.Close()
	// the os.Stdout would be an io.Writer used to write uncompressed data to
	io.Copy(os.Stdout, r)

	*/

}
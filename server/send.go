package main

import (
	"net"
	"fmt"
	"bufio"
	"io"
	"os"
)

func dial(ip string, port string, file io.Reader){
	var b [1000]byte
	nb := 0
	buf := bufio.NewReader(file)
	for {
		c, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		if err == nil {
			b[nb] = c
			nb++
		}
	}
	message := b[0:nb]
	conn, _ := net.Dial("tcp", ip + ":" + port)
	conn.Write(message)
	status, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}

func main(){
	f, _ := os.Open("send.go")
	defer f.Close()
	dial("127.0.0.1", "4376", f)
}
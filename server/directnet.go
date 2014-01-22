package main
 
import (
	"net"
	"os"
)
 
const (
	RECV_BUF_LEN = 1024
)

func main() {
 
	listener, err := net.Listen("tcp", "0.0.0.0:4376")
	if err != nil {
		println("error listening:", err.Error())
		os.Exit(1)
	}
 
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		go EchoFunc(conn)
	}
}
 
func EchoFunc(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	n, err := conn.Read(buf)
	if err != nil {
		println("Error reading:", err.Error())
		return
	}
	println("received ", n, " bytes of data =", string(buf))
 
	//send reply
	_, err = conn.Write(buf)
	if err != nil {
		println("Error send reply:", err.Error())
	}else {
		println("Reply sent")
	}
}
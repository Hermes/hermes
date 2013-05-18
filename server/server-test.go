package main 

import (
	"fmt"
	"net"
	"os"
)

func getIPs(hostname string) {
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(err)
		//return "", ""
	}

	localIP := fmt.Sprintf("%s", addrs[0])
	fmt.Println(localIP)


}

func main() {
	name, _ := os.Hostname();
	fmt.Println(name)
	getIPs(name)
}
package server

import (
	"fmt"
	"github.com/secondbit/wendy"
	"os"
	"net"
)

type Server struct {
	Cluster wendy.Cluster
	Node wendy.Node
	Hostname string
	ID string
	localIP, globalIP string
}

func NewServer() Server {
	Server := Server{}
	Server.Hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	Server.ID, err = wendy.NodeIDFromBytes([]byte(hostname))
	if err != nil {
		panic(err)
	}

	Server.localIP, Server.globalIP = getIPs()
	Server.Node = wendy.NewNode(Server.ID, Server.localIP,
								Server.globalIP, "angelhack",
								31337)

	credentials := NewCredentials(Server.Id)

	Server.Cluster = wendy.NewCluster(node, credentials)
	

}

func getIPs(hostname string) (string, string) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Println(err)
		return "",""
	}

	localIP := fmt.Sprintf("%s", addrs[0])
	globalIP := "38.116.199.162"

	return localIP, globalIP
}
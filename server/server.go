package server

import (
	"fmt"
	"github.com/secondbit/wendy"
	"os"
	"net"
)

var (
	Server HermesServer
	FileMap Files
)

func init() {
	Server = NewServer()
	FileMap = NewFilesMap()
}

type HermesServer struct {
	Cluster wendy.Cluster
	Node wendy.Node
	Hostname string
	ID string
	localIP, globalIP string
	kill chan bool
}

func NewServer() HermesServer {
	Server := HermesServer{}
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
								1337)

	credentials := wendy.Passphrase("Hermes")
	Server.Cluster = wendy.NewCluster(node, credentials)

	go func() {
		defer cluster.Stop()
		err := cluster.Listen()
		if err != nil {
			panic(err.Error())
		}
	}()

	app := &HermesApplication{}
	Server.Cluster.RegisterCallback(app)
	Server.Cluster.Join("ip of another Node", 1337)

	return Server
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
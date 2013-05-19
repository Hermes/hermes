package server 

import (
	"fmt"
	"github.com/secondbit/wendy"
)

type HermesApplication struct {
}

func (app *HermesApplication) OnError(err error) {
    panic(err.Error())
}

func (app *HermesApplication) OnDeliver(msg wendy.Message) {
    fmt.Println("Received message: ", msg)
}

func (app *HermesApplication) OnForward(msg *wendy.Message, next wendy.NodeID) bool {
    fmt.Printf("Forwarding message %s to Node %s.", msg.ID, next)
    return true // return false if you don't want the message forwarded
}

func (app *HermesApplication) OnNewLeaves(leaves []*wendy.Node) {
    fmt.Println("Leaf set changed: ", leaves)
}

func (app *HermesApplication) OnNodeJoin(node *wendy.Node) {
    fmt.Println("Node joined: ", node.ID)
}

func (app *HermesApplication) OnNodeExit(node *wendy.Node) {
    fmt.Println("Node left: ", node.ID)
}

func (app *HermesApplication) OnHeartbeat(node *wendy.Node) {
    fmt.Println("Received heartbeat from ", node.ID)
}

app := &HermesApplication{}
cluster.RegisterCallback(app)
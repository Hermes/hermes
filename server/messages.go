package server 

import (
	//"fmt"
	"github.com/secondbit/wendy"
	"encoding/json"
)

/** Wendy Message Legend **/
// Purpose 16 - PushMsg 
// Purpose 17 - PullMsg 
// Purpose 18 - ManifestRequest
// Purpose 19 - ManifestResponse

const (
	PUSH_TX 		=	byte(16)
	PUSH_RX			=	byte(17)
	PULL_TX	 		=	byte(18)
	PULL_RX			=	byte(19)
	MANIFEST_TX 	=	byte(20)
	MANIFEST_RX		=	byte(21)
)

type FileMessage struct {
	Status string
	VaultID string
	Filename string
	SplitID string
	RedundantLevel string
	Data []byte
}

type ManifestMessage struct {
	VaultID string
	Filenames []string
}

/****** MESSAGE ROUTER/HANDLERS ******/

func MessageHandler(msg wendy.Message) {
	var response wendy.Message

	switch msg.Purpose {
		case PUSH_TX:
			respData := HandlePUSH_TX(msg.Value)
			response = Server.Cluster.NewMessage(PUSH_RX, msg.Sender.ID)

		case PUSH_RX:
			//Handle PUSH response
		case PULL_TX:
			//Handle PULL message
		case PULL_RX:
			//Handler PULL response
		case MANIFEST_TX:
			//Handle Manifest request
		case MANIFEST_RX:
			//Handle Manifest response
	}

	Server.Cluster.send(response)
}

func HandlePUSH_TX(msg []byte) []byte {
	var push FileMessage
	err := json.Unmarshal(msg, &push)
	if err != nil {
		panic(err)
	}

	err = FileMap.Push(push)
	resp := FileMessage{}
	if err != nil {
		resp.Status = "Failed to push file"
	}

	resp.Status = "Sucessfully pushed file"
	resp.VaultID = push.VaultID
	resp.Filename = push.Filename
	resp.SplitID = push.SplitID
	resp.RedundantLevel = push.RedundantLevel
	
	respJSON, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return respJSON

}
func HandlePUSH_RX(msg []byte) {}


func Push(VaultID string, filepath string) {
	
}

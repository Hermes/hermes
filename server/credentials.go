package server

import (
	"encoding/hex"
	"crypto/sha256"
	"fmt"
	"time"
	"os"
)

type NetCredentials struct {
	NetID []byte
}

func (f *NetCredentials) SetID(ID []byte) {
    f.NetID = ID
}

func NewCredentials() NetCredentials {
	hostname, _ := os.Hostname()
	micro := time.Now().Nanosecond()

	VaultIDString := fmt.Sprintf("%s%s", hostname, micro)
	sha := sha256.New()
	sha.Write([]byte(VaultIDString))
	VaultID := sha.Sum(nil)

	cred := NetCredentials{}
    cred.SetID(VaultID)
	return cred
}

func (c NetCredentials) Valid(NetID []byte) bool {
	/*if VaultID.length == 32 {
		return true
	}
	return false*/
	return true
}

func (c NetCredentials) Marshal() []byte {
	return c.NetID
}

func (c NetCredentials) String() string {
	return hex.EncodeToString(c.NetID)
}

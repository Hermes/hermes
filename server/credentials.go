package server 

import (
	"fmt"
	"crypto/sha256"
	"time"
	"crypto/hex"
)

type NetCredentials struct {
	NetID []byte
}

func NewCredentials() NetCredentials {
	hostname, err := os.Hostname();
	micro := time.Now().Nanosecond()

	VaultIDString := fmt.Sprintf("%s%s", hostname, micro)
	sha := sha256.New()
	sha.Write([]byte(VaultIDString))
	VaultID := sha.Sum()

	return Credentials(VauldID)
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


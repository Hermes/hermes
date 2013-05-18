package server 

import (
	"fmt"
	"crypto/sha256"
	"time"
	"crypto/hex"
)

type Credentials struct {
	VaultID []byte
}

func NewCredentials() Credentials {
	hostname, err := os.Hostname();
	micro := time.Now().Nanosecond()

	VaultIDString := fmt.Sprintf("%s%s", hostname, micro)
	sha := sha256.New()
	sha.Write([]byte(VaultIDString))
	VaultID := sha.Sum()

	return Credentials(VauldID)
}

func (c Credentials) Valid(VaultID []byte) bool {
	if VaultID.length == 32 {
		return true
	}
	return false
}

func (c Credentials) Marshal() []byte {
	return c.VaultID
}

func (c Credentials) String() string {
	return hex.EncodeToString(c.VaultID)
}


package names

import (
	"crypto/rand"
	"fmt"
)

func GetId() string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	return fmt.Sprintf("%x", randomBytes)
}

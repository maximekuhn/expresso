package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func generateSessionId() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

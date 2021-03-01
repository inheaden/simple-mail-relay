package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GetRandomString returns a rondom string of length length.
func GetRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// ValidateNonce validates if the payload (nonce + secret) matches the provided hash
func ValidateNonce(nonce string, hash string, secret string) bool {
	payload := nonce + secret
	generatedHash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(generatedHash[:]) == hash
}

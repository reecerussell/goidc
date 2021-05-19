package util

import (
	"crypto/sha256"
	"encoding/base64"
)

// Sha256 returns a SHA-256 hash of value, represented in base64.
func Sha256(value string) string {
	alg := sha256.New()
	alg.Write([]byte(value))

	hash := alg.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}

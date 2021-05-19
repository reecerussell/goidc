package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256(t *testing.T) {
	const testString = "Hello World"
	const expectedHash = "pZGm1Av0IEBKARczz7exkNYsZb8LzaMrV7J32a2fFG4="

	hash := Sha256(testString)
	assert.Equal(t, expectedHash, hash)
}

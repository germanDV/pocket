package pocket

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

// SafeCompare performs a constant-time comparison of two strings to protect against timing attacks.
// It hashes both strings to ensure they have the same length.
func SafeCompare(token1, token2 string) bool {
	h1 := sha256.Sum256([]byte(token1))
	h2 := sha256.Sum256([]byte(token2))
	return subtle.ConstantTimeCompare(h1[:], h2[:]) == 1
}

// GenerateString generates a random string of the specified length.
// If for any reason `rand.Read` fails, this function will panic!
func GenerateString(len int) string {
	bytes := make([]byte, len)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

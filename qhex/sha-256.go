package qhex

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(input, yan string) string {
	hash := sha256.New()
	// Write the input string to the hash object
	input += yan
	hash.Write([]byte(input))

	// Get the resulting hash as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

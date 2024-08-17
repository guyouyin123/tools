package aes

import (
	"encoding/base64"
	"encoding/hex"
)

type CipherText []byte

const (
	PrintHex = iota
	PrintBase64
)

func (ct CipherText) hexEncode() string {
	return hex.EncodeToString(ct)
}

func (ct CipherText) base64Encode() string {
	return base64.StdEncoding.EncodeToString(ct)
}

func hexDecode(cipherText string) ([]byte, error) {
	return hex.DecodeString(cipherText)
}

func base64Decode(cipherText string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cipherText)
}

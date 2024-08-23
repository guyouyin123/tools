package conf

import (
	"encoding/base64"
	"encoding/hex"
)

type CipherText []byte

const (
	PrintHex = iota
	PrintBase64
)

func (ct CipherText) HexEncode() string {
	return hex.EncodeToString(ct)
}

func (ct CipherText) Base64Encode() string {
	return base64.StdEncoding.EncodeToString(ct)
}

func HexDecode(cipherText string) ([]byte, error) {
	return hex.DecodeString(cipherText)
}

func Base64Decode(cipherText string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cipherText)
}

package des

import "bytes"

const (
	CBCMode = iota
	//CFBMode
	CTRMode
	ECBMode
	//OFBMode
)

type FillMode int

const (
	PkcsZero FillMode = iota
	Pkcs7
	Pkcs5
)

// The blockSize argument should be 16, 24, or 32.
// Corresponding AES-128, AES-192, or AES-256.
func (fm FillMode) pkcs7Padding(plainText []byte, blockSize int) []byte {
	paddingSize := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(plainText, paddingText...)
}

func (fm FillMode) unPkcs7Padding(plainText []byte) []byte {
	length := len(plainText)
	number := int(plainText[length-1])
	return plainText[:length-number]
}

func (fm FillMode) zeroPadding(plainText []byte, blockSize int) []byte {
	if plainText[len(plainText)-1] == 0 {
		return nil
	}
	paddingSize := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(0)}, paddingSize)
	return append(plainText, paddingText...)
}

func (fm FillMode) unZeroPadding(plainText []byte) []byte {
	length := len(plainText)
	count := 1
	for i := length - 1; i > 0; i-- {
		if plainText[i] == 0 && plainText[i-1] == plainText[i] {
			count++
		}
	}
	return plainText[:length-count]
}

// PKCS5Padding填充模式
func (fm FillMode) pkcs5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := make([]byte, len(plainText)+padding)
	copy(padText, plainText)
	for i := len(plainText); i < len(padText); i++ {
		padText[i] = byte(padding)
	}
	return padText
}

// PKCS5Padding去除填充模式
func (fm FillMode) unPkcs5Padding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:length-unpadding]
}

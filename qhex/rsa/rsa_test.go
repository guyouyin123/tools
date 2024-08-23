package goencrypt

import (
	"fmt"
	"github.com/guyouyin123/tools/qhex/conf"
	"testing"
)

const defaultPublicFile = "./testdata/pub.txt"
const defaultPrivateFile = "./testdata/pri.txt"

func TestNewRSACipher(t *testing.T) {
	// rsa encryption
	cipher := NewRSACipher(conf.PrintBase64, defaultPublicFile, defaultPrivateFile)
	cipherText, err := cipher.RSAEncrypt([]byte("hello world"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cipherText)

	// rsa decryption
	plainText, err := cipher.RSADecrypt(cipherText)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(plainText)
}

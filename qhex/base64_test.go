package qhex

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestPageTokenEncode(t *testing.T) {
	token := PageTokenEncode(10, 20)
	fmt.Println(token)
	offset, page, pageSize := PageTokenDecode(token)
	fmt.Println(offset, page, pageSize)

}

func TestAes(t *testing.T) {
	data, err := Encrypt("123456", "e10adc3949ba59abbe56e057f20f883e")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}

func TestAes2(t *testing.T) {
	key := []byte("e10adc3949ba59abbe56e057f20f883e") // 16字节的密钥
	plaintext := []byte("123456")

	ciphertext, err := AESCBCEncrypt(plaintext, key)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	fmt.Println("Ciphertext (hex):", hex.EncodeToString(ciphertext))
}

func TestMD5(t *testing.T) {
	a := MD52("123456", "aa")
	fmt.Println(a)
}

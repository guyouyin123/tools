package des

import (
	"fmt"
	"testing"
)

func TestNewDESCipher(t *testing.T) {
	key := []byte("0123456789abcdef") //密码
	iv := []byte("0123456789abcdef")  //偏移量
	model := ECBMode                  //模式
	pkcs := Pkcs7                     //填充
	out := PrintHex                   //输出

	cipher := NewDESCipher(key, iv, model, pkcs, out)
	cipherText, err := cipher.DESEncrypt([]byte("hello world"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cipherText)

	// des decryption
	plainText, err := cipher.DESDecrypt(cipherText)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(plainText)
}

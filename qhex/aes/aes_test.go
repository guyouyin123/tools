package aes

import (
	"fmt"
	"github.com/guyouyin123/tools/qhex/conf"
	"testing"
)

func TestCipherAES_AESEncrypt(t *testing.T) {
	// aes decryption
	key := []byte("123456789abcdefg") //密码
	iv := []byte("0123456789abcdef")  //偏移量
	model := conf.CTRMode             //模式
	pkcs := conf.Pkcs5                //填充
	out := conf.PrintHex              //输出
	cipher, err := NewAESCipher(key, iv, model, pkcs, out)
	if err != nil {
		fmt.Println(err)
		return
	}
	cipherText, err := cipher.AESEncrypt([]byte("hello world"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cipherText)

	// aes decryption
	plainText, err := cipher.AESDecrypt(cipherText)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(plainText)
}

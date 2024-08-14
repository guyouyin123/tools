package qhex

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func Encrypt(data, key string) (result string, err error) {
	var resultByte []byte
	resultByte, err = aesEncrypt([]byte(data), []byte(key))
	if err != nil {
		return
	}
	result = base64.StdEncoding.EncodeToString(resultByte)
	return
}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	//origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AESCBCEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 生成一个随机的初始化向量(IV)
	iv := []byte("1234567890123456")
	fmt.Println("iv:", string(iv))

	// 创建一个加密的CBC模式
	mode := cipher.NewCBCEncrypter(block, iv)

	// 对明文进行填充
	plaintext = PKCS7Padding(plaintext, aes.BlockSize)

	// 加密明文
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 返回初始化向量(IV)和密文
	return append(iv, ciphertext...), nil
}

// PKCS7Padding 实现PKCS7填充
func PKCS7Padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

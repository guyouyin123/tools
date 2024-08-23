package aes

import (
	"crypto/aes"
	"errors"
	"github.com/guyouyin123/tools/qhex/conf"
	"strconv"
)

type CipherAES struct {
	conf.Cipher
}

const BlockSize = 16

func NewAESCipher(key, iv []byte, groupMode int, fillMode conf.FillMode, decodeType int) (*CipherAES, error) {
	if len(iv) != BlockSize {
		return nil, errors.New("IV length must equal block size")
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("crypto/aes: invalid key size " + strconv.Itoa(len(key)))
	}
	return &CipherAES{conf.Cipher{
		GroupMode:  groupMode,
		FillMode:   fillMode,
		DecodeType: decodeType,
		Key:        key,
		Iv:         iv,
	},
	}, nil
}

// AESEncrypt aes加密
func (c *CipherAES) AESEncrypt(plainText []byte) (cipherText string, err error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return
	}
	plainData := c.Fill(plainText, block.BlockSize())
	if plainData == nil {
		err = errors.New("unsupported content to be encrypted")
		return
	}
	c.Output = make(conf.CipherText, len(plainData))
	if err = c.Encrypt(block, plainData); err != nil {
		return
	}
	return c.Encode(), nil
}

// AESDecrypt aes解密
func (c *CipherAES) AESDecrypt(cipherText string) (plainText string, err error) {
	cipherData, err := c.Decode(cipherText)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return
	}
	if len(cipherData)%block.BlockSize() != 0 {
		err = errors.New("cipher text is not a multiple of the block size")
		return
	}
	if err = c.Decrypt(block, cipherData); err != nil {
		return
	}
	plainData, err := c.UnFill(c.Output)
	if err != nil {
		return "", conf.HandleError(err)
	}
	return string(plainData), nil
}

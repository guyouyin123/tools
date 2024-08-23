package goencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/guyouyin123/tools/qhex/conf"
	"os"
)

type CipherRSA struct {
	PubKey string
	PriKey string
	conf.Cipher
}

// New returns a CipherRSA pointer, the keyFile parameters is address of
// public key or private key.
// If the length of keyFile is 1, which means RSA encryption.
// If the length of keyFile greater than 1, which means RSA decryption.
// Meanwhile, the second element of keyFile represents private key.
func NewRSACipher(decodeType int, keyFile ...string) *CipherRSA {
	if len(keyFile) > 1 {
		return &CipherRSA{
			keyFile[0],
			keyFile[1],
			conf.Cipher{
				DecodeType: decodeType,
			},
		}
	} else {
		return &CipherRSA{
			keyFile[0],
			"",
			conf.Cipher{
				DecodeType: decodeType,
			},
		}
	}
}

// Encrypt encrypts the given message with RSA and the padding
// scheme from PKCS#1 v1.5.
func (cr *CipherRSA) RSAEncrypt(plainText []byte) (cipherText string, err error) {
	data, err := os.ReadFile(cr.PubKey)
	if err != nil {
		return
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	cr.Output, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), plainText)
	if err != nil {
		return
	}
	return cr.Encode(), nil
}

// Decrypt parses the given message with RSA private key in PKCS#1, ASN.1 DER form.
func (cr *CipherRSA) RSADecrypt(cipherText string) (plainText string, err error) {
	data, err := os.ReadFile(cr.PriKey)
	if err != nil {
		return
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return "", errors.New("private key error")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	cipherData, err := cr.Decode(cipherText)
	if err != nil {
		return
	}
	plainData, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipherData)
	if err != nil {
		return
	}
	return string(plainData), nil
}

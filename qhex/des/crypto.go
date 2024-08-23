package des

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"strings"
)

type Crypto interface {
	Encrypt(plainText []byte) (string, error)
	Decrypt(cipherText string) (string, error)
}

type Cipher struct {
	GroupMode  int
	FillMode   FillMode
	DecodeType int
	Key        []byte
	Iv         []byte
	Output     CipherText
}

func (c *Cipher) Encrypt(block cipher.Block, plainData []byte) (err error) {
	c.Output = make([]byte, len(plainData))
	if c.GroupMode == CBCMode {
		cipher.NewCBCEncrypter(block, c.Iv).CryptBlocks(c.Output, plainData)
		return
	}
	if c.GroupMode == ECBMode {
		c.NewECBEncrypter(block, plainData)
		return
	}
	if c.GroupMode == CTRMode {
		c.NewCTREncrypter(block, plainData)
		return
	}
	return
}

func (c *Cipher) Decrypt(block cipher.Block, cipherData []byte) (err error) {
	c.Output = make([]byte, len(cipherData))
	if c.GroupMode == CBCMode {
		cipher.NewCBCDecrypter(block, c.Iv).CryptBlocks(c.Output, cipherData)
		return
	}
	if c.GroupMode == ECBMode {
		c.NewECBDecrypter(block, cipherData)
		return
	}
	if c.GroupMode == CTRMode {
		c.NewCTRDecrypter(block, cipherData)
		return
	}
	return
}

// Encode
// default print format is base64
func (c *Cipher) Encode() string {
	if c.DecodeType == PrintHex {
		return c.Output.hexEncode()
	} else {
		return c.Output.base64Encode()
	}
}

func (c *Cipher) Decode(cipherText string) ([]byte, error) {
	if c.DecodeType == PrintBase64 {
		return base64Decode(cipherText)
	} else if c.DecodeType == PrintHex {
		return hexDecode(cipherText)
	} else {
		return nil, errors.New("unsupported print type")
	}
}

func (c *Cipher) Fill(plainText []byte, blockSize int) []byte {
	//填充模式
	switch c.FillMode {
	case PkcsZero:
		return c.FillMode.zeroPadding(plainText, blockSize)
	case Pkcs5:
		return c.FillMode.pkcs5Padding(plainText, blockSize)
	case Pkcs7:
		return c.FillMode.pkcs7Padding(plainText, blockSize)
	default:
		return nil
	}
}

func (c *Cipher) UnFill(plainText []byte) (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	switch c.FillMode {
	case PkcsZero:
		return c.FillMode.unZeroPadding(plainText), nil
	case Pkcs5:
		return c.FillMode.unPkcs5Padding(plainText), nil
	case Pkcs7:
		return c.FillMode.unPkcs7Padding(plainText), nil

	default:
		return nil, errors.New("unsupported fill mode")
	}
}

func (c *Cipher) NewECBEncrypter(block cipher.Block, plainData []byte) {
	tempText := c.Output
	for len(plainData) > 0 {
		block.Encrypt(tempText, plainData[:block.BlockSize()])
		plainData = plainData[block.BlockSize():]
		tempText = tempText[block.BlockSize():]
	}
}

func (c *Cipher) NewECBDecrypter(block cipher.Block, cipherData []byte) {
	tempText := c.Output
	for len(cipherData) > 0 {
		block.Decrypt(tempText, cipherData[:block.BlockSize()])
		cipherData = cipherData[block.BlockSize():]
		tempText = tempText[block.BlockSize():]
	}
}

func (c *Cipher) NewCTREncrypter(block cipher.Block, plaintext []byte) {
	if len(plaintext) > 0 {
		stream := cipher.NewCTR(block, c.Iv)
		ciphertext := make([]byte, len(plaintext))
		stream.XORKeyStream(ciphertext, plaintext)
		c.Output = ciphertext
	}
}

func (c *Cipher) NewCTRDecrypter(block cipher.Block, cipherData []byte) {
	if len(cipherData) > 0 {
		stream := cipher.NewCTR(block, c.Iv)
		plaintext := make([]byte, len(cipherData))
		stream.XORKeyStream(plaintext, cipherData)
		c.Output = plaintext
	}
}

func (c *Cipher) NewCFBEncrypter(block cipher.Block, plaintext []byte) {
	if len(plaintext) > 0 {
		stream := cipher.NewCFBEncrypter(block, c.Iv)
		ciphertext := make([]byte, len(plaintext))
		stream.XORKeyStream(ciphertext, plaintext)
		c.Output = ciphertext
	}
}

const runtimeErr = "runtime error:"

func handleError(err error) error {
	if strings.HasPrefix(err.Error(), runtimeErr) {
		return errors.New("encrypted and decrypted passwords are inconsistent")
	}
	return err
}

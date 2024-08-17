package qhex

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1(input, yan string) string {
	o := sha1.New()
	input += yan
	o.Write([]byte(input))
	return hex.EncodeToString(o.Sum(nil))
}

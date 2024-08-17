package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"sync"
)

var md5Pool = sync.Pool{
	New: func() interface{} {
		return md5.New()
	},
}

func MD5(str, yan string) string {
	/*
		yan:盐
		使用对象池
	*/
	h := md5Pool.Get().(hash.Hash)
	str = str + yan
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
func MD5Sum(str, yan string) string {
	/*
		一次性计算给定数据 data 的 MD5 值,不需要手动维护中间状态
	*/
	str = str + yan
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

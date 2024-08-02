package qhex

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// Base64Encode base64编码
func Base64Encode(s string) string {
	encodedMessage := base64.StdEncoding.EncodeToString([]byte(s))
	return encodedMessage
}

// Base64Decode base64解码
func Base64Decode(s string) (string, error) {
	decodedMessage, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decodedMessage), nil
}

// 编码流式翻页token
func PageTokenEncode(page, pageSize int) string {
	p := PageToken{
		Page:      page,
		PageSize:  pageSize,
		Timestamp: time.Now().Unix(),
	}
	jsonStr, _ := json.Marshal(p)
	s := Base64Encode(string(jsonStr))
	return s
}

// 解码流式翻页token--size 为每页默认大小
func PageTokenDecode(pageToken string) (int, int, int) {
	if pageToken == "" {
		return 0, 0, 20
	}
	s, err := Base64Decode(pageToken)
	if err != nil {
		return 0, 0, 20
	}
	var p PageToken
	_ = json.Unmarshal([]byte(s), &p)

	offset := (p.Page - 1) * p.PageSize
	return offset, p.Page, p.PageSize
}

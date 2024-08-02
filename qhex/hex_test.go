package qhex

import (
	"fmt"
	"testing"
)

func TestPageTokenEncode(t *testing.T) {
	token := PageTokenEncode(10, 20)
	fmt.Println(token)
	offset, page, pageSize := PageTokenDecode(token)
	fmt.Println(offset, page, pageSize)

}

func TestDemo(t *testing.T) {

}

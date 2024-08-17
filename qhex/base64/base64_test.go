package base64

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

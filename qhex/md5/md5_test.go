package md5

import (
	"fmt"
	"github.com/guyouyin123/tools/qhex/sha"
	"testing"
)

func TestMD5(t *testing.T) {
	a := MD5Sum("123456", "aa")
	fmt.Println(a)

	b := sha.SHA1("123456", "aa")
	fmt.Println(b)

	c := sha.SHA256("123456", "aa")
	fmt.Println(c)
}

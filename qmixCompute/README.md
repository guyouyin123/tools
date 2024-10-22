# 公式计算相关工具,做到了防止精度丢失问题

```go
import (
	"fmt"
	"testing"
	qmixCompute "github.com/guyouyin123/tools/qmixCompute"
)

func TestMixCompute(t *testing.T) {
	result := qmixCompute.MixCompute("a*b+c", map[rune]float64{
		'a': float64(100),
		'b': 0.5,
		'c': float64(10),
	})
	fmt.Println(result)

	result2 := qmixCompute.MixCompute("a/b", map[rune]float64{
		'a': float64(10),
		'b': 3,
	})
	fmt.Println(result2)
}
```


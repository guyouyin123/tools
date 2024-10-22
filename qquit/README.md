# 优雅推出程序

```go
import (
	"fmt"
	"testing"
	"time"
)

func TestQuitSignal(t *testing.T) {
	QuitSignal(demo)
}

func demo() {
	fmt.Println("程序退出做的事情")
	time.Sleep(time.Second * 3)
}
```



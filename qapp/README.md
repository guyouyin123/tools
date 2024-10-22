# 市场app应用信息获取的工具

```go
import (
	"fmt"
	"github.com/guyouyin123/tools/qapp"
	"testing"
)

func TestCheckUrl(t *testing.T) {
	url := "https://apps.apple.com/app/id1559181149"
	id, source, err := qapp.CheckUrl(url)
	if err != nil {
		fmt.Println(err)
	}
	switch source {
	case "en", "cn":
		info, _ := qapp.GetItunesAppleUrl(id, source)
		fmt.Println(info.Results[0]) //{晶核 [游戏 角色扮演 动作] https://is1-ssl.mzstatic.com/image/thumb/Purple221/v4/ca/02/a1/ca02a113-4ed9-ed15-c1d2-e05f45d41e11/AppIcon-1x_U007emarketing-0-7-0-85-220-0.png/512x512bb.jpg com.hermes.p6game}
	case "google":
		googleData, _ := qapp.GetGoogleUrl(id)
		fmt.Println(*googleData)
	}
}
```
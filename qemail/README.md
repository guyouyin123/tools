**安装**

```go
go get github.com/guyouyin123/tools
tools.qemail
使用的：
go get github.com/jordan-wright/email
```

我们需要额外一些工作。我们知道邮箱使用`SMTP/POP3/IMAP`等协议从邮件服务器上拉取邮件。邮件并不是直接发送到邮箱的，而是邮箱请求拉取的。 所以，我们需要配置`SMTP/POP3/IMAP`服务器。从头搭建固然可行，而且也有现成的开源库，但是比较麻烦。现在一般的邮箱服务商都开放了`SMTP/POP3/IMAP`服务器。 我这里拿 126 邮箱来举例，使用`SMTP`服务器。当然，用 QQ 邮箱也可以。

- 首先，登录邮箱；
- 点开顶部的设置，选择`POP3/SMTP/IMAP`；
- 点击开启`IMAP/SMTP`服务，按照步骤开启即可，有个密码设置，记住这个密码，后面有用。


**开始开发**
```go
package main

import (
	"fmt"
	"github.com/guyouyin123/tools/qemail"
)

func main() {
	fromEmail := "xxoo@163.com"
	host := "smtp.163.com"
	pwd := "123"

	e := qemail.Init(fromEmail, host, pwd)

	toS := []string{
		"xxoo@qq.com",
	}
	title := "邮件测试主题名称"
	content := []byte("邮件内容")
	attachFile := "./user.xlsx" //附件
	err := e.Send(toS, title, content, nil, []string{"xxoo@163.com"}, nil,attachFile)
	if err != nil {
		fmt.Println(err)
	}
}

```

**demo**
```go
func TestEmail_Send(t *testing.T) {
	fromEmail := "xxoo@163.com"
	host := "smtp.163.com"
	//host := "imap.163.com"
	pwd := "pwd"

	e := Init(fromEmail, host, pwd)

	toS := []string{
		"1009122179@qq.com",
	}
	title := "邮件测试主题名称"
	content := []byte("邮件内容")
	err := e.Send(toS, title, content, nil, []string{"guyouyin@163.com"}, nil,"")
	if err != nil {
		fmt.Println(err)
	}
}
```

**demo2**

```go
func TestEmail_html(t *testing.T) {
	fromEmail := "xxoo@163.com"
	host := "smtp.163.com"
	//host := "imap.163.com"
	pwd := "pwd"

	e := Init(fromEmail, host, pwd)

	toS := []string{
		"xxoo@qq.com",
		"xxoo@163.com",
	}
	title := "邮件测试主题名称"
	contentHtml := []byte(`<li><a "https://www.cnblogs.com/guyouyin123">Jeff技术栈</a></li>`)
	err := e.Send(toS, title, nil, contentHtml, []string{"xxoo@163.com"}, nil,"")
	if err != nil {
		fmt.Println(err)
	}
}
```


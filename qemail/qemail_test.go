package qemail

import (
	"fmt"
	"testing"
)

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
	attachFile := "./user.xlsx" //附件
	err := e.Send(toS, title, content, nil, []string{"xxoo@163.com"}, nil, attachFile)
	if err != nil {
		fmt.Println(err)
	}
}

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
	err := e.Send(toS, title, nil, contentHtml, []string{"xxoo@163.com"}, nil, "")
	if err != nil {
		fmt.Println(err)
	}
}

package qemail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type Email struct {
	FromEmail string //发放方
	Host      string //host地址
	Pwd       string //授权码
	Addr      string
	E         *email.Email
	Auth      smtp.Auth
}

func Init(FromEmail, Host, Pwd string) *Email {
	e := email.NewEmail()
	e.From = fmt.Sprintf("<%s>", FromEmail)
	auth := smtp.PlainAuth("", FromEmail, Pwd, Host)
	return &Email{
		FromEmail: FromEmail,
		Host:      Host,
		Pwd:       Pwd,
		Addr:      fmt.Sprintf("%s:25", Host),
		E:         e,
		Auth:      auth,
	}
}

/*
Send 发送邮件
@param toS 接收方
@param title 标题
@param content 内容
@param CC 抄送
@param BCC 秘密抄送
@param attachFile 附件路径
*/
func (this *Email) Send(toS []string, title string, content, contentHtml []byte, CC, BCC []string, attachFile string) error {
	e := this.E
	if len(CC) > 0 {
		e.Cc = CC
	}
	if len(BCC) > 0 {
		e.Bcc = BCC
	}
	e.To = toS
	e.Subject = title
	if len(content) > 0 {
		e.Text = content
	}
	if len(contentHtml) > 0 {
		e.HTML = contentHtml
	}
	if attachFile != "" {
		_, err := e.AttachFile(attachFile)
		if err != nil {
			return err
		}
	}
	err := e.Send(this.Addr, this.Auth)
	if err != nil {
		return err
	}
	return nil
}

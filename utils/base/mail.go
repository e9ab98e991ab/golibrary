/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//
package base

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
)

// 邮件信道
type MailChan struct {
	ToEmail         string
	Title           string
	Body            string
	Rebackinterfack RebackMail
}

type RebackMail interface {
	BackMail(toEmail, title, body string, b bool, err error) error
}

type Mail struct {
	Host     string
	Port     string
	Email    string
	Password string
	Sender   string
}

func (m *Mail) Send(toEmail, title, body string, backInterface RebackMail) {
	func() {
		if err := recover(); err != nil {
			errstr := "golang Mail::send throw a fatal error"
			fmt.Println(errstr)
			backInterface.BackMail(toEmail, title, body, false, errors.New(errstr))
			return
		}
	}()
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	host := m.Host
	email := m.Email
	password := m.Password
	from := mail.Address{m.Sender, email}
	to := mail.Address{"接收人", toEmail}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(title)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)
	fmt.Println("send email...")
	err := smtp.SendMail(
		host+":"+m.Port,
		auth,
		email,
		[]string{to.Address},
		[]byte(message),
	)
	fmt.Println("send email over!")
	if err != nil {
		backInterface.BackMail(toEmail, title, body, false, err)
	} else {
		backInterface.BackMail(toEmail, title, body, true, nil)
	}
	return
}

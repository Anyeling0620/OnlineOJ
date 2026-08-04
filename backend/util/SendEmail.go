package util

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
)

// SendCode 发送验证码
func SendCode(toUserEmail string, code string) error {
	e := email.NewEmail()
	e.From = "Test <" + os.Getenv("EMAIL_USERNAME") + ">"
	e.To = []string{toUserEmail}
	e.Subject = "OnlineOJ"
	e.HTML = []byte("您的验证码是<h1>" + code + " </h1>")
	return e.SendWithTLS(os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_PORT"), smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")), &tls.Config{InsecureSkipVerify: true, ServerName: os.Getenv("EMAIL_HOST")})
}

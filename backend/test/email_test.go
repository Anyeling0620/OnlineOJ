package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

var (
	EMAIL_USERNAME = ""
	EMAIL_PASSWORD = ""
	EMAIL_HOST     = ""
	EMAIL_PORT     = ""

	EMAIL_RECEIVE_TEST = ""
)

//
//func TestSendEmail(t *testing.T) {
//	e := email.NewEmail()
//	e.From = "Test <" + os.Getenv("EMAIL_USERNAME") + ">"
//	e.To = []string{os.Getenv("EMAIL_RECEIVE_TEST")}
//	e.Subject = "验证码发送测试"
//	e.HTML = []byte("您的验证码是<h1> 114514 </h1>")
//	err := e.SendWithTLS(os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_PORT"), smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")), &tls.Config{InsecureSkipVerify: true, ServerName: os.Getenv("EMAIL_HOST")})
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Test <" + EMAIL_USERNAME + ">"
	e.To = []string{EMAIL_RECEIVE_TEST}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码是<h1> 114514 </h1>")
	err := e.SendWithTLS(EMAIL_HOST+":"+EMAIL_PORT, smtp.PlainAuth("", EMAIL_USERNAME, EMAIL_PASSWORD, EMAIL_HOST), &tls.Config{InsecureSkipVerify: true, ServerName: EMAIL_HOST})
	if err != nil {
		t.Fatal(err)
	}
}

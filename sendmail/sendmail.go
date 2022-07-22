package sendmail

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type SendMail struct {
	dialer *gomail.Dialer
	host   string
	port   int
	user   string
	passwd string
}

const (
	ContentTypeHtml = "text/html"
	ContentTypeText = "text"
)

func NewGMail(host string, port int, user, passwd string) *SendMail {

	sendmail := &SendMail{
		host:   host,
		port:   port,
		user:   user,
		passwd: passwd,
	}
	sendmail.conn()

	return sendmail
}

func (p SendMail) conn() {
	p.dialer = gomail.NewDialer(p.host, p.port, p.user, p.passwd)
	p.dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}

func (p *SendMail) Mail(from string, to, cc []string, subject, content, contentType string) error {
	p.conn()

	msger := gomail.NewMessage()
	msger.SetHeader("From", from)
	msger.SetHeader("To", to...)

	if len(cc) > 0 {
		for _, addr := range cc {
			msger.SetAddressHeader("Cc", addr, addr)
		}
	}
	msger.SetHeader("Subject", subject)

	msger.SetBody(contentType, content)

	sc, err := p.dialer.Dial()
	if err != nil {
		return err
	}
	defer sc.Close()

	if err = gomail.Send(sc, msger); err != nil {
		return err
	}

	return nil
}

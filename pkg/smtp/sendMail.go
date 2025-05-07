package smtp

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
)

type MailService struct {
	email    string
	auth     smtp.Auth
	smtpHost string
	smtpPort string
}

func NewMailService(username, password string) *MailService {
	return &MailService{
		auth:     smtp.PlainAuth("", username, password, "smtp.gmail.com"),
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		email:    username,
	}
}

func (s *MailService) sendMail(to []string, bodyBuf bytes.Buffer) error {

	// Sending email.
	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, s.auth, s.email, to, bodyBuf.Bytes())
	return err
}

func (s *MailService) SendMail(to []string, subject, text string) error {

	var bodyBuf bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	bodyBuf.Write(fmt.Appendf(nil, "Subject: %s \n%s\n\n", subject, mimeHeaders))

	// write raw text to bodyBuf
	bodyBuf.Write([]byte(text))

	err := s.sendMail(to, bodyBuf)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent successfully")
	}

	return err
}

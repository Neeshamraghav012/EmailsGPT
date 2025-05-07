package mail

import (
	"log"
	"project/models"
	"project/pkg/smtp"
)

type MailService struct {
	client  *smtp.MailService
	to      []string
	subject string
}

func NewMailService(username, password string, to []string) *MailService {
	return &MailService{
		client:  smtp.NewMailService(username, password),
		to:      to,
		subject: "Critical Emails Summary",
	}
}

func (s *MailService) Notify(req models.SummarizerResponse) error {
	if req.Emails == nil {
		log.Println("Received nil request for Slack notification service")
		return nil
	}

	formattedMessage := s.buildMailMessage(req)
	if formattedMessage == "" {
		log.Println("Failed to format message for Slack")
		return nil
	}

	err := s.client.SendMail(s.to, s.subject, formattedMessage)
	if err != nil {
		log.Println("Failed to send email:", err)
	}

	return nil
}

func (s *MailService) buildMailMessage(response models.SummarizerResponse) string {
	message := "<H2>Hey there! Here are some critical emails:</H2><hr/>"
	for _, email := range response.Emails {
		message += "Subject: " + email.Subject + "<br/>"
		message += "From: " + email.From + "<br/>"
		message += "Summary: " + email.Summary + "<br/>"
		message += "Link: " + email.Link + "`<br/><br/>"
		message += "<hr/>"
	}

	return message
}

package gmail

import (
	"log"
	"project/models"
	googleimap "project/pkg/google_imap"
	"time"
)

type Service struct {
	client *googleimap.Client
}

func NewService(client *googleimap.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FetchEmails(since time.Time) (mails []models.Mail, err error) {
	logTag := "GmailService"
	if since.IsZero() {
		log.Println(logTag + "Received zero time for fetching emails")
		return nil, nil
	}
	pkgMails, err := s.client.FetchMailsSince(since)
	if err != nil {
		log.Printf(logTag+"Error fetching emails: %v", err)
		return
	}

	if len(pkgMails) == 0 {
		log.Println(logTag + "No emails found")
		return
	}

	mails = getSvcMails(pkgMails)

	return
}

func getSvcMails(mails []googleimap.Mail) []models.Mail {
	svcMails := make([]models.Mail, 0)
	for _, mail := range mails {
		svcMails = append(svcMails, models.Mail{
			Subject:  mail.Subject,
			From:     mail.From,
			BodyText: mail.BodyText,
			Link:     mail.Link,
		})
	}
	return svcMails
}

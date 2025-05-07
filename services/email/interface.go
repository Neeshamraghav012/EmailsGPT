package email

import (
	"project/models"
	"time"
)

type Interface interface {
	FetchEmails(since time.Time) ([]models.Mail, error)
}

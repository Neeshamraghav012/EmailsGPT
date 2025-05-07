package summarizer

import "project/models"

type Request struct {
	Emails   []models.Mail
	Criteria []models.Criteria
}

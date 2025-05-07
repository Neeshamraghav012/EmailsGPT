package notifier

import "project/models"

type Interface interface {
	Notify(models.SummarizerResponse) error
}


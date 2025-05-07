package handler

import (
	"context"
	"log"
	"project/models"
	"project/services/email"
	"project/services/notifier"
	"project/services/summarizer"
	"time"
)

type CronHandler struct {
	duration          time.Duration
	criteria          []models.Criteria
	emailService      email.Interface
	summarizerService summarizer.Interface
	notifierService   notifier.Interface
}

func NewCronHandler(
	duration time.Duration,
	criteria []models.Criteria,
	emailService email.Interface,
	summarizerService summarizer.Interface,
	notifierService notifier.Interface,
) *CronHandler {
	return &CronHandler{
		duration:          duration + 10*time.Second,
		criteria:          criteria,
		emailService:      emailService,
		summarizerService: summarizerService,
		notifierService:   notifierService,
	}
}

func (h *CronHandler) Process(ctx context.Context) error {
	logTag := "Function: CronHandler :: Process :: "

	// 1. Fetch emails using the email service
	emails, err := h.emailService.FetchEmails(time.Now().Add(-h.duration))
	if err != nil {
		log.Printf(logTag+"Error fetching emails: %v", err)
		return err
	}

	if len(emails) == 0 {
		log.Println(logTag + "No emails to process")
		return nil
	}

	results, err := h.summarizerService.Summarize(ctx, &summarizer.Request{
		Emails:   emails,
		Criteria: h.criteria,
	})
	if err != nil {
		log.Printf(logTag+"Error summarizing emails: %v", err)
		return err
	}

	if len(results.Emails) == 0 {
		log.Println(logTag + "No results to notify")
		return nil
	}

	// 2. Notify using the notifier service
	err = h.notifierService.Notify(results)
	if err != nil {
		log.Printf(logTag+"Error notifying: %v", err)
		return err
	}

	return nil
}

func (h *CronHandler) Start(ctx context.Context, shutdownChannel chan struct{}) {
	logTag := "Function: CronHandler :: Start :: "
	ticker := time.NewTicker(h.duration)
	defer ticker.Stop()

	// First run
	if err := h.Process(ctx); err != nil {
		log.Printf(logTag+"Error processing cron job: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			log.Println(logTag + "Processing cron job")
			if err := h.Process(ctx); err != nil {
				log.Printf(logTag+"Error processing cron job: %v", err)
			}
		case <-ctx.Done():
			log.Println(logTag + "Stopping cron job")
			shutdownChannel <- struct{}{}
			return
		}
	}

}

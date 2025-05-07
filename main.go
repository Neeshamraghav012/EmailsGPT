package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/config"
	"project/handler"
	"project/models"
	googleimap "project/pkg/google_imap"
	"project/services/email/gmail"
	mailNotifier "project/services/notifier/mail"
	"project/services/summarizer/openai"

	"log"
)

func main() {
	// Initialize the processor
	duration := time.Duration(3 * time.Hour)

	// Define the criteria for processing
	criteria := []models.Criteria{
		{
			Name:        "Social",
			Description: "Linkedin invitations",
		},
	}

	processor := initializeProcessor(criteria, duration)

	ctx, cancel := context.WithCancel(context.Background())
	shutdownChannel := make(chan struct{})
	defer close(shutdownChannel)
	// Start the processor
	go processor.Start(ctx, shutdownChannel)

	// listen to interrupt signal and close the ctx

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
	log.Println("Received interrupt signal, shutting down...")
	cancel()
	log.Println("Shutdown complete")

}

func initializeProcessor(criteria []models.Criteria, duration time.Duration) handler.Interface {
	// Initialize the email service
	emailFetchService := gmail.NewService(
		googleimap.NewClient(
			config.GetEnv("GMAIL_USER"),
			config.GetEnv("GMAIL_PASSWORD"),
		),
	)

	// Initialize the summarizer service
	summarizerService := openai.NewOpenAIService()

	// Initialize the notifier service
	notifierService := mailNotifier.NewMailService(
		config.GetEnv("GMAIL_USER"),
		config.GetEnv("GMAIL_PASSWORD"),
		[]string{config.GetEnv("NOTIFICATION_EMAIL")},
	)

	// Create the cron handler
	processor := handler.NewCronHandler(
		duration,
		criteria,
		emailFetchService,
		summarizerService,
		notifierService,
	)

	return processor
}

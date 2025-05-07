package slack

import (
	"log"
	"project/models"
	"project/pkg/slack"
)

type SlackService struct {
	client *slack.SlackClient
}

func NewSlackService(client *slack.SlackClient) *SlackService {
	return &SlackService{
		client: client,
	}
}

func (s *SlackService) Notify(req models.SummarizerResponse) error {
	if req.Emails == nil {
		log.Println("Received nil request for Slack notification service")
		return nil
	}

	formattedMessage := s.buildSlackMessage(req)
	if formattedMessage == "" {
		log.Println("Failed to format message for Slack")
		return nil
	}

	s.client.SendMessage(formattedMessage)
	return nil
}

func (s *SlackService) buildSlackMessage(response models.SummarizerResponse) string {
	// TODO: Implement the logic to build a well formatted Slack message from the response
	message := "Hey there! Here are some critical emails:\n\n"
	for _, email := range response.Emails {
		message += "Subject: `" + email.Subject + "`\n"
		message += "From: `" + email.From + "`\n"
		message += "Summary: `" + email.Summary + "`\n"
		message += "Link: `" + email.Link + "`\n\n"
	}

	return message
}

package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SlackClient struct {
	webhookURL string
}

func NewSlackClient(webhookURL string) *SlackClient {
	return &SlackClient{
		webhookURL: webhookURL,
	}
}

func (c *SlackClient) SendMessage(message string) (err error) {
	// Implement the logic to send a message to Slack using the webhook URL

	body := SlackMessage{
		Text: message,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error sending request to Slack:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Slack webhook failed: %s\n", resp.Status)
		data := make([]byte, 1024)
		_, err := resp.Body.Read(data)
		if err != nil {
			log.Println("Error reading response body:", err)
		}
		return fmt.Errorf("slack webhook failed: %+v , Status: %d", data, resp.StatusCode)
	}

	log.Println("Message successfully sent to Slack channel!")
	return nil
}

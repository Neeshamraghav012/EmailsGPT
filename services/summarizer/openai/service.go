package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"project/config"
	"project/models"
	"project/services/summarizer"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService() *OpenAIService {
	openAIClient := openai.NewClient(
		option.WithAPIKey(config.GetEnv("OPENAI_API_KEY")),
	)

	return &OpenAIService{
		client: &openAIClient,
	}
}

func (s *OpenAIService) Summarize(
	ctx context.Context,
	req *summarizer.Request,
) (resp models.SummarizerResponse, err error) {
	resp.Emails = make([]models.EmailResponseEntry, 0)
	prompt := s.getPrompt(req)

	if prompt == "" {
		log.Println("Received empty prompt for OpenAI summarization")
		return
	}

	response, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		log.Printf("Error generating summary: %v", err)
		return
	}
	if response == nil {
		log.Println("Received nil response from OpenAI")
		return
	}

	resp, err = s.processResponse(req, response)
	if err != nil {
		log.Printf("Error processing OpenAI response: %v", err)
		return
	}

	return
}

func (s *OpenAIService) getPrompt(req *summarizer.Request) (prompt string) {
	// TODO : Implement the logic to build a prompt for OpenAI
	prompt = `I am providing you subjects of emails and their id number allotted to them.
Emails:
%s
Please mention the emails that are critical as per given the following criteria. 
%s

Please just output the json in the following schema:
{
	"critical_emails": [
		{
			"id": 1,
			"subject": "Subject of email",
			"from": "From email address",
			"summary": "mention why it is critical"
		}
	]
}

Don't mention anything else. Just output the json. Without even using the word "json" and markdown.
`
	if len(req.Emails) == 0 {
		log.Println("No emails provided for summarization")
		return ""
	}
	if len(req.Criteria) == 0 {
		log.Println("No criteria provided for summarization")
		return ""
	}

	emails := ""
	for i, email := range req.Emails {
		emails += fmt.Sprintf("\tEmail %d: %s\n", i+1, email.Subject)
		emails += fmt.Sprintf("\t\t(From: %s)\n", email.From)
	}

	criteria := ""
	for i, crit := range req.Criteria {
		criteria += fmt.Sprintf("%d. %s: %s\n", i+1, crit.Name, crit.Description)
	}
	prompt = fmt.Sprintf(prompt, emails, criteria)

	return
}

func (s *OpenAIService) processResponse(req *summarizer.Request, response *openai.ChatCompletion) (resp models.SummarizerResponse, err error) {
	if response == nil {
		log.Println("Received nil response from OpenAI")
		return
	}

	if len(response.Choices) == 0 {
		log.Println("No choices found in OpenAI response")
		return
	}

	if len(response.Choices[0].Message.Content) == 0 {
		log.Println("Empty content in OpenAI response")
		return
	}

	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &resp)
	if err != nil {
		log.Printf("Error unmarshalling OpenAI response: %v", err)
		return
	}

	idToLinkMap := make(map[int]string)
	for i, email := range req.Emails {
		idToLinkMap[i+1] = email.Link
	}

	for i, email := range resp.Emails {
		if link, exists := idToLinkMap[email.ID]; exists {
			resp.Emails[i].Link = link
		} else {
			resp.Emails[i].Link = ""
		}
	}

	return
}

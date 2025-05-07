package summarizer

import (
	"context"
	"project/models"
)

type Interface interface {
	Summarize(ctx context.Context, req *Request) (models.SummarizerResponse, error)
}

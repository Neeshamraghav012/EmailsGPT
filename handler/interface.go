package handler

import (
	"context"
)

type Interface interface {
	Process(ctx context.Context) error
	Start(ctx context.Context, stop chan struct{})
}

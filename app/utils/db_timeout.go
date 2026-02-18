package utils

import (
	"context"
	"time"
)

const QueryTimeout = 20 * time.Second

func WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, QueryTimeout)
}

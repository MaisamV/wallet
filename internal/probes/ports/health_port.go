package ports

import (
	"context"
	"time"
)

// DatabaseChecker defines the interface for checking database connectivity
type DatabaseChecker interface {
	CheckDatabase(ctx context.Context) (bool, time.Duration, error)
}

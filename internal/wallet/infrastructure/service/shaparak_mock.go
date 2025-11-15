package service

import (
	"context"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofrs/uuid/v5"
	"math/rand"
	"time"
)

type ShaparakMockService struct {
	logger logger.Logger
}

func NewShaparakMockService(logger logger.Logger) *ShaparakMockService {
	return &ShaparakMockService{
		logger: logger,
	}
}

func (s *ShaparakMockService) Withdraw(ctx context.Context, userId int64, idempotency *uuid.UUID, withdrawAmount int64) (*uuid.UUID, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return mockHttpCall(opCtx, userId, idempotency, withdrawAmount)
}

func mockHttpCall(ctx context.Context, id int64, idempotency *uuid.UUID, amount int64) (*uuid.UUID, error) {
	resultCh := make(chan uuid.UUID)
	go func() {
		if rand.Float64() < 0.20 {
			//timeout
			time.Sleep(5001 * time.Millisecond)
		} else {
			//success
			time.Sleep(50 * time.Millisecond)
			u, err := uuid.NewV7()
			if err == nil {
				resultCh <- u
			}
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		return &result, nil
	}
}

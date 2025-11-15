package service

import (
	"context"
	"github.com/gofrs/uuid/v5"
)

type BankService interface {
	Withdraw(ctx context.Context, userId int64, idempotency *uuid.UUID, withdrawAmount int64) (*uuid.UUID, error)
}

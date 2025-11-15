package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/MaisamV/wallet/internal/wallet/ports/repo"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofrs/uuid/v5"
	"time"
)

type WithdrawCommand struct {
	UserId      int64
	Amount      int64
	Idempotency *uuid.UUID
	ReleaseTime *time.Time
}

func (cc *WithdrawCommand) Err() error {
	if cc.Amount <= 0 {
		return errors.New("amount cannot be negative or zero")
	}
	if cc.Idempotency == nil {
		return errors.New("idempotency cannot be null")
	}
	if cc.ReleaseTime == nil || cc.ReleaseTime.Before(time.Now()) {
		return errors.New("release time must not be in the past")
	}
	return nil
}

type WithdrawCommandHandler struct {
	logger logger.Logger
	repo   repo.WalletWriter
}

func NewWithdrawCommandHandler(logger logger.Logger, repo repo.WalletWriter) *WithdrawCommandHandler {
	return &WithdrawCommandHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *WithdrawCommandHandler) Handle(ctx context.Context, command WithdrawCommand) (*uuid.UUID, error) {
	if err := command.Err(); err != nil {
		return nil, fmt.Errorf("input variables are not correct: %w", err)
	}
	txnID, err := h.repo.Debit(ctx, command.UserId, command.Idempotency, command.Amount, command.ReleaseTime)
	if err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}
	return txnID, nil
}

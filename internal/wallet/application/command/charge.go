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

type ChargeCommand struct {
	UserId      int64
	Amount      int64
	Idempotency *uuid.UUID
	ReleaseTime *time.Time
}

func (cc *ChargeCommand) IsValid() bool {
	return cc.Amount > 0 && cc.Idempotency != nil && (cc.ReleaseTime == nil || cc.ReleaseTime.After(time.Now()))
}

func (cc *ChargeCommand) Err() error {
	if cc.Amount > 0 {
		return errors.New("amount cannot be negative or zero")
	}
	if cc.Idempotency != nil {
		return errors.New("idempotency cannot be null")
	}
	if cc.ReleaseTime == nil || cc.ReleaseTime.After(time.Now()) {
		return errors.New("release time must not be in the past")
	}
	return nil
}

type ChargeCommandHandler struct {
	logger logger.Logger
	repo   repo.WalletWriter
}

func NewAddContactCommandHandler(logger logger.Logger, repo repo.WalletWriter) *ChargeCommandHandler {
	return &ChargeCommandHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *ChargeCommandHandler) Handle(ctx context.Context, command ChargeCommand) (*uuid.UUID, error) {
	if err := command.Err(); err != nil {
		return nil, fmt.Errorf("input variables are not correct: %w", err)
	}
	txnID, err := h.repo.Charge(ctx, command.UserId, command.Idempotency, command.Amount, command.ReleaseTime)
	if err != nil {
		return nil, fmt.Errorf("failed to charge wallet: %w", err)
	}
	return txnID, nil
}

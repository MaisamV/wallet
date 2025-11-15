package query

import (
	"context"
	"fmt"
	"github.com/MaisamV/wallet/internal/wallet/entity"
	"github.com/MaisamV/wallet/internal/wallet/ports/repo"
	"github.com/MaisamV/wallet/platform/logger"
)

type GetBalanceQuery struct {
	UserID int64
}

type GetBalanceQueryHandler struct {
	logger logger.Logger
	repo   repo.WalletReader
}

func NewGetBalanceQueryHandler(logger logger.Logger, repo repo.WalletReader) *GetBalanceQueryHandler {
	return &GetBalanceQueryHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *GetBalanceQueryHandler) Handle(ctx context.Context, query GetBalanceQuery) (*entity.Wallet, error) {
	wallet, err := h.repo.GetBalance(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return wallet, nil
}

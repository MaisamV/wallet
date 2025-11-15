package query

import (
	"context"
	"fmt"
	"github.com/MaisamV/wallet/internal/wallet/entity"
	"github.com/MaisamV/wallet/internal/wallet/ports/repo"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofrs/uuid/v5"
)

type GetTransactionPageQuery struct {
	UserID int64
	Cursor *uuid.UUID
	Limit  int
}

type GetTransactionPageQueryHandler struct {
	logger logger.Logger
	repo   repo.WalletReader
}

func NewGetTransactionPageQueryHandler(logger logger.Logger, repo repo.WalletReader) *GetTransactionPageQueryHandler {
	return &GetTransactionPageQueryHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *GetTransactionPageQueryHandler) Handle(ctx context.Context, query GetTransactionPageQuery) (*entity.TransactionPage, error) {
	page, err := h.repo.GetTransactionList(ctx, query.UserID, query.Cursor, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction page: %w", err)
	}
	return page, nil
}

package command

import (
	"context"
	"github.com/MaisamV/wallet/internal/wallet/entity"
	"github.com/MaisamV/wallet/internal/wallet/ports/repo"
	"github.com/MaisamV/wallet/platform/logger"
)

type ReleaseCommand struct {
	BatchSize int
}

type ReleaseCommandHandler struct {
	logger logger.Logger
	repo   repo.WalletWriter
}

func NewReleaseCommandHandler(logger logger.Logger, repo repo.WalletWriter) *ReleaseCommandHandler {
	return &ReleaseCommandHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *ReleaseCommandHandler) Handle(ctx context.Context, command ReleaseCommand) ([]entity.Transaction, error) {
	releasedTxns, err := h.repo.ReleaseDueTransactions(ctx, command.BatchSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("release failed")
		return nil, err
	}

	if len(releasedTxns) > 0 {
		h.logger.Info().Int("count", len(releasedTxns)).Msg("released transactions")
	}

	return releasedTxns, nil
}

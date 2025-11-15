package command

import (
	"context"
	"fmt"
	"github.com/MaisamV/wallet/internal/wallet/entity"
	"github.com/MaisamV/wallet/internal/wallet/ports/repo"
	"github.com/MaisamV/wallet/internal/wallet/ports/service"
	"github.com/MaisamV/wallet/platform/logger"
)

type WithdrawCommand struct {
	Limit int
}

type WithdrawCommandHandler struct {
	logger      logger.Logger
	repo        repo.WalletRepo
	bankService service.BankService
	workerCount int
	pendingCh   chan *entity.Transaction
}

func NewWithdrawCommandHandler(logger logger.Logger, repo repo.WalletRepo, bankService service.BankService, workerCount int) *WithdrawCommandHandler {
	return &WithdrawCommandHandler{
		logger:      logger,
		repo:        repo,
		workerCount: workerCount,
		bankService: bankService,
		pendingCh:   make(chan *entity.Transaction),
	}
}

func (h *WithdrawCommandHandler) Start() {
	for i := 0; i < h.workerCount; i++ {
		go h.WorkerLoop()
	}
}

func (h *WithdrawCommandHandler) Handle(ctx context.Context, command WithdrawCommand) error {
	pendingTxs, err := h.repo.GetPendingTransactions(ctx, command.Limit)
	if err != nil {
		return fmt.Errorf("failed to read pending transactions: %w", err)
	}
	for _, tx := range pendingTxs {
		h.pendingCh <- &tx
	}
	return nil
}

func (h *WithdrawCommandHandler) WorkerLoop() {
	ctx := context.Background()
	for tx := range h.pendingCh {
		bankTxUUID, err := h.bankService.Withdraw(ctx, tx.UserID, &tx.Idempotency, tx.Amount)
		if err != nil {
			h.logger.Error().Err(err).Msg("error happened while trying to call Bank API")
			tx.RetryCount++
			if tx.RetryCount >= 5 {
				err := h.repo.UpdateTransactionStatus(ctx, &tx.ID, entity.FAILED, nil)
				if err != nil {
					h.logger.Error().Err(err).Msg("couldn't update transaction status to failed")
				}
			} else {
				err := h.repo.IncreaseTransactionRetryCount(ctx, &tx.ID)
				if err != nil {
					h.logger.Error().Err(err).Msg("couldn't update transaction status to failed")
				}
			}
		} else {
			err := h.repo.UpdateTransactionStatus(ctx, &tx.ID, entity.SUCCESS, bankTxUUID)
			if err != nil {
				h.logger.Error().Err(err).Msg("couldn't update transaction status to success")
			}
			h.logger.Info().Str("id", tx.ID.String()).Msg("successfully withdraw")
		}
	}
}

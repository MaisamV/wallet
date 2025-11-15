package user

import (
	"github.com/MaisamV/wallet/internal/wallet/application/command"
	"github.com/MaisamV/wallet/internal/wallet/application/query"
	infrastructure "github.com/MaisamV/wallet/internal/wallet/infrastructure/repo"
	service2 "github.com/MaisamV/wallet/internal/wallet/infrastructure/service"
	"github.com/MaisamV/wallet/internal/wallet/presentation/http"
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProvideWalletRepository(logger logger.Logger, db *pgxpool.Pool) *infrastructure.PgxWalletRepo {
	return infrastructure.NewPgxWalletRepo(logger, db)
}

func ProvideShaparakMockService(logger logger.Logger) *service2.ShaparakMockService {
	return service2.NewShaparakMockService(logger)
}

func ProvideChargeCommandHandler(logger logger.Logger, repo *infrastructure.PgxWalletRepo) *command.ChargeCommandHandler {
	return command.NewChargeCommandHandler(logger, repo)
}

func ProvideDebitCommandHandler(logger logger.Logger, repo *infrastructure.PgxWalletRepo) *command.DebitCommandHandler {
	return command.NewDebitCommandHandler(logger, repo)
}

func ProvideReleaseCommandHandler(logger logger.Logger, repo *infrastructure.PgxWalletRepo) *command.ReleaseCommandHandler {
	return command.NewReleaseCommandHandler(logger, repo)
}

func ProvideWithdrawCommandHandler(logger logger.Logger, repo *infrastructure.PgxWalletRepo, service *service2.ShaparakMockService, cfg *config.Config) *command.WithdrawCommandHandler {
	return command.NewWithdrawCommandHandler(logger, repo, service, cfg.Withdraw.WorkerCount)
}

func ProvideGetBalanceQueryHandler(logger logger.Logger, repo *infrastructure.PgxWalletRepo) *query.GetBalanceQueryHandler {
	return query.NewGetBalanceQueryHandler(logger, repo)
}

func ProvideGetTransactionPageQueryHandler(logger logger.Logger, repo *infrastructure.PgxWalletRepo) *query.GetTransactionPageQueryHandler {
	return query.NewGetTransactionPageQueryHandler(logger, repo)
}

func ProvideWalletHandler(logger logger.Logger, withdrawHandler *command.DebitCommandHandler,
	chargeHandler *command.ChargeCommandHandler, balanceHandler *query.GetBalanceQueryHandler,
	transactionPageHandler *query.GetTransactionPageQueryHandler) *http.WalletHandler {
	return http.NewWalletHandler(logger, withdrawHandler, chargeHandler, balanceHandler, transactionPageHandler)
}

// WalletSet is a wire provider set for all user dependencies
var WalletSet = wire.NewSet(
	ProvideWalletRepository,
	ProvideDebitCommandHandler,
	ProvideChargeCommandHandler,
	ProvideReleaseCommandHandler,
	ProvideWithdrawCommandHandler,
	ProvideShaparakMockService,
	ProvideGetBalanceQueryHandler,
	ProvideGetTransactionPageQueryHandler,
	ProvideWalletHandler,
)

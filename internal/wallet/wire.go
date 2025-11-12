package user

import (
	infrastructure "github.com/MaisamV/wallet/internal/wallet/infrastructure/repo"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProvideWalletRepository(logger logger.Logger, db *pgxpool.Pool) *infrastructure.PgxWalletRepo {
	return infrastructure.NewPgxWalletRepo(logger, db)
}

// WalletSet is a wire provider set for all user dependencies
var WalletSet = wire.NewSet(
	ProvideWalletRepository,
)

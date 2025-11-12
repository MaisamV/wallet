package repo

import (
	"context"
	"github.com/MaisamV/wallet/internal/wallet/entity"
	"github.com/gofrs/uuid/v5"
	"time"
)

type WalletRepo interface {
	WalletWriter
	WalletReader
}

type WalletWriter interface {
	Charge(ctx context.Context, userId int64, idempotency *uuid.UUID, chargeAmount int64, releaseTime *time.Time) (txnId *uuid.UUID, err error)
	Debit(ctx context.Context, userId int64, debitAmount int64, releaseTime time.Time) error
}

type WalletReader interface {
	GetWallet(ctx context.Context, userId uuid.UUID) (*entity.Wallet, error)
	GetTransactionList(ctx context.Context, userId uuid.UUID) ([]entity.Transaction, error)
}

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
	Debit(ctx context.Context, userId int64, idempotency *uuid.UUID, debitAmount int64, releaseTime *time.Time) (txnId *uuid.UUID, err error)
}

type WalletReader interface {
	GetBalance(ctx context.Context, userId int64) (*entity.Wallet, error)
	GetTransactionList(ctx context.Context, userId int64, cursor *uuid.UUID, limit int) (*entity.TransactionPage, error)
}

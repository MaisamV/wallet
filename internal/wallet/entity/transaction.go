package entity

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Transaction struct {
	ID          uuid.UUID
	WalletID    int64
	UserID      int64
	Type        string
	Amount      int64
	Idempotency uuid.UUID
	ReleaseTime time.Time
	Released    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

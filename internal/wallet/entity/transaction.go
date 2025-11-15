package entity

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type TransactionType = string

const (
	CREDIT TransactionType = "credit"
	DEBIT                  = "debit"
)

type Status = string

const (
	pending  Status = "pending"
	CANCELED        = "canceled"
	FAILED          = "failed"
	SUCCESS         = "success"
)

type Transaction struct {
	ID          uuid.UUID
	WalletID    int64
	UserID      int64
	Type        TransactionType
	Status      Status
	Amount      int64
	Idempotency uuid.UUID
	ReleaseTime time.Time
	Released    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TransactionPage struct {
	TransactionList []Transaction
	Cursor          *uuid.UUID
}

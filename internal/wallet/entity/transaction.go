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
	ID          uuid.UUID       `json:"id,omitempty"`
	WalletID    int64           `json:"-"`
	UserID      int64           `json:"-"`
	Type        TransactionType `json:"type,omitempty"`
	Status      Status          `json:"status,omitempty"`
	Amount      int64           `json:"amount,omitempty"`
	Idempotency uuid.UUID       `json:"-"`
	ReleaseTime *time.Time      `json:"release_time,omitempty"`
	Released    bool            `json:"released"`
	RetryCount  int
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"-"`
}

type TransactionPage struct {
	TransactionList []Transaction `json:"transaction_list,omitempty"`
	Cursor          *uuid.UUID    `json:"cursor,omitempty"`
}

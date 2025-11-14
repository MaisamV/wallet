package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/MaisamV/wallet/internal/wallet/entity"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PgxWalletRepo struct {
	logger logger.Logger
	db     *pgxpool.Pool
}

func NewPgxWalletRepo(logger logger.Logger, db *pgxpool.Pool) *PgxWalletRepo {
	return &PgxWalletRepo{
		logger: logger,
		db:     db,
	}
}

// Charge Adds credit to the user wallet
func (dc *PgxWalletRepo) Charge(ctx context.Context, userId int64, idempotency *uuid.UUID, chargeAmount int64, releaseTime *time.Time) (*uuid.UUID, error) {
	opCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if idempotency == nil {
		return nil, errors.New("charge operations must have idempotency")
	}

	if chargeAmount <= 0 {
		return nil, errors.New("negative or 0 is not acceptable amount for charge operation")
	}

	if releaseTime != nil && time.Now().After(*releaseTime) {
		return nil, errors.New("release time can't be in the past")
	}

	var query string
	if releaseTime == nil {
		query = chargeQuery
	} else {
		query = chargeWithReleaseQuery
	}
	var transactionID uuid.UUID
	err := dc.db.QueryRow(opCtx, query, userId, chargeAmount, releaseTime, idempotency).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("database charge operation failed: %w", err)
	}

	return &transactionID, nil
}

// Debit deducts from user's wallet balance
func (dc *PgxWalletRepo) Debit(ctx context.Context, userId int64, idempotency *uuid.UUID, debitAmount int64, releaseTime *time.Time) (*uuid.UUID, error) {
	opCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if idempotency == nil {
		return nil, errors.New("debits operations must have idempotency")
	}

	if debitAmount <= 0 {
		return nil, errors.New("negative or 0 is not acceptable amount for debit operation")
	}

	if releaseTime == nil {
		return nil, errors.New("debits must have release time")
	}

	if time.Now().After(*releaseTime) {
		return nil, errors.New("release time can't be in the past")
	}

	var transactionID uuid.UUID
	err := dc.db.QueryRow(opCtx, debitWithReleaseQuery, userId, debitAmount, releaseTime, idempotency).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("database debit operation failed: %w", err)
	}

	return &transactionID, nil
}

// GetBalance return user's wallet balance
func (dc *PgxWalletRepo) GetBalance(ctx context.Context, userId int64) (*entity.Wallet, error) {
	opCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var w *entity.Wallet
	var id int64
	var totalBalance int64
	var availableBalance int64
	err := dc.db.QueryRow(opCtx, getBalance, userId).Scan(&id, &userId, &totalBalance, &availableBalance)
	switch err {
	case pgx.ErrNoRows:
		w = entity.NewWallet(int64(0), userId, int64(0), int64(0))
	case nil:
		w = entity.NewWallet(id, userId, totalBalance, availableBalance)
	default:
		return nil, fmt.Errorf("get balance operation failed: %w", err)
	}

	return w, nil
}

// GetTransactionList return a list of user transactions
func (dc *PgxWalletRepo) GetTransactionList(ctx context.Context, userId int64, cursor *uuid.UUID, limit int) (*entity.TransactionPage, error) {
	opCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if limit > 30 {
		limit = 30
	}

	var err error
	var rows pgx.Rows
	if cursor == nil {
		rows, err = dc.db.Query(opCtx, getTransactionsFirstPage, userId, limit)
	} else {
		rows, err = dc.db.Query(opCtx, getTransactionsNextPage, userId, limit, cursor)
	}
	if err != nil {
		return nil, fmt.Errorf("get transaction list operation failed: %w", err)
	}
	defer rows.Close()
	list := make([]entity.Transaction, 0, limit)
	for rows.Next() {
		t := entity.Transaction{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &t.Status, &t.Amount, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("error in reading transaction row: %w", err)
		}
		list = append(list, t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("something went wrong reading transaction list: %w", rows.Err())
	}

	page := entity.TransactionPage{
		TransactionList: list,
	}
	size := len(list)
	if size == limit {
		page.Cursor = &list[size-1].ID
	}

	return &page, nil
}

// Close gracefully close all database pool connections
func (dc *PgxWalletRepo) Close() {
	dc.db.Close()
}

const (
	chargeQuery = `
WITH upserted_wallet AS (
    INSERT INTO wallets (user_id, total_balance, available_balance)
    VALUES ($1, $2, $2)
    ON CONFLICT (user_id) DO UPDATE
    SET total_balance = wallets.total_balance + EXCLUDED.total_balance,
        available_balance = wallets.available_balance + EXCLUDED.total_balance,
        updated_at = NOW()
    RETURNING id AS wallet_id, user_id
),
inserted_txn AS (
    INSERT INTO transactions 
        (wallet_id, user_id, type, status, amount, release_time, released, idempotency_key)
    SELECT wallet_id, user_id, 'credit' AS type, 'success' AS status, $2 AS amount, $3 AS release_time, TRUE, $4 AS idempotency_key
    FROM upserted_wallet 
    RETURNING id AS txn_id
)
SELECT txn_id FROM inserted_txn;
`
	chargeWithReleaseQuery = `
WITH upserted_wallet AS (
    INSERT INTO wallets (user_id, total_balance, available_balance)
    VALUES ($1, $2, 0)
    ON CONFLICT (user_id) DO UPDATE
    SET total_balance = wallets.total_balance + EXCLUDED.total_balance,
        updated_at = NOW()
    RETURNING id AS wallet_id, user_id
),
inserted_txn AS (
    INSERT INTO transactions 
        (wallet_id, user_id, type, status, amount, release_time, released, idempotency_key)
    SELECT wallet_id, user_id, 'credit' AS type, 'blocked' AS status, $2 AS amount, $3 AS release_time, FALSE, $4 AS idempotency_key
    FROM upserted_wallet 
    RETURNING id AS txn_id
)
SELECT txn_id FROM inserted_txn;
`
	debitWithReleaseQuery = `
WITH updated_wallet AS (
    UPDATE wallets
    SET
        available_balance = available_balance - $2,
        updated_at = NOW()
    WHERE user_id = $1 AND available_balance >= $2
    RETURNING id AS wallet_id, user_id
),
inserted_txn AS (
    INSERT INTO transactions 
        (wallet_id, user_id, type, status, amount, release_time, released, idempotency_key)
    SELECT wallet_id, user_id, 'debit' AS type, 'blocked' AS status, ($2 * -1) AS amount, $3 AS release_time, FALSE, $4 AS idempotency_key
    FROM updated_wallet
    RETURNING id AS txn_id
)
SELECT txn_id FROM inserted_txn;
`
	getBalance = `
SELECT id, user_id, total_balance, available_balance 
FROM wallets
WHERE user_id = $1
`
	getTransactionsFirstPage = `
SELECT id, user_id, type, status, amount, created_at 
FROM transactions
WHERE user_id = $1
ORDER BY ID DESC
LIMIT $2
`
	getTransactionsNextPage = `
SELECT id, user_id, type, status, amount, created_at 
FROM transactions
WHERE user_id = $1
AND id < $3
ORDER BY ID DESC
LIMIT $2
`
)

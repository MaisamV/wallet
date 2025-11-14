package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofrs/uuid/v5"
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

	query := debitWithReleaseQuery
	var transactionID uuid.UUID
	err := dc.db.QueryRow(opCtx, query, userId, debitAmount, releaseTime, idempotency).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("database debit operation failed: %w", err)
	}

	return &transactionID, nil
}

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
)

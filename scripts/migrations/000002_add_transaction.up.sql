BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    user_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('credit', 'debit')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('blocked', 'failed', 'cancelled', 'success')),
    amount             bigint NOT NULL,
    release_time       TIMESTAMPTZ NULL,
    released           BOOLEAN NOT NULL DEFAULT FALSE,
    idempotency_key    UUID NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_txn_user_created ON transactions (user_id, created_at DESC);
CREATE INDEX idx_txn_release_pending ON transactions (release_time, released)
    WHERE released = FALSE;
CREATE UNIQUE INDEX idx_txn_user_key ON transactions (user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;


COMMIT;
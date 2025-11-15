BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    user_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('credit', 'debit')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'failed', 'success')),
    retry_count        int NOT NULL DEFAULT 0,
    last_retry         TIMESTAMPTZ NULL,
    amount             bigint NOT NULL,
    release_time       TIMESTAMPTZ NULL,
    released           BOOLEAN NOT NULL DEFAULT FALSE,
    idempotency_key    UUID NOT NULL,
    bank_response_id   UUID NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_transactions_pending_retry ON transactions (status, last_retry, id);
CREATE INDEX idx_transactions_user_id_id ON transactions (user_id, id DESC);
CREATE INDEX idx_txn_user_created ON transactions (user_id, created_at DESC);
CREATE INDEX idx_txn_release_pending ON transactions (release_time, released)
    WHERE released = FALSE;
CREATE UNIQUE INDEX idx_txn_user_key ON transactions (user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;


COMMIT;
BEGIN;

-- The uuid-ossp extension is required for the uuid_generate_v4 function.
--CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    total_balance bigint NOT NULL DEFAULT 0,
    available_balance bigint NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_wallet_user_id ON wallets(user_id);

COMMIT;
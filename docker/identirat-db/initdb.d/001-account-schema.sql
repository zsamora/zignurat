CREATE TABLE IF NOT EXISTS accounts (
    id          BIGSERIAL PRIMARY KEY,
    username    TEXT NOT NULL,
    password    TEXT NOT NULL,
    module      SMALLINT NOT NULL DEFAULT 0,
    acc_role    SMALLINT NOT NULL DEFAULT 0,
    owner_uuid  UUID,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    jwt_token   TEXT,
    refr_token  TEXT,

    CONSTRAINT chk_module   CHECK (module BETWEEN 0 AND 255),
    CONSTRAINT chk_acc_role CHECK (acc_role BETWEEN 0 AND 255)
    );

CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts (deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_username
    ON accounts (username) WHERE deleted_at IS NULL;
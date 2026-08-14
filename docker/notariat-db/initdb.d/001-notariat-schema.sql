CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    names       TEXT NOT NULL,
    surnames    TEXT NOT NULL,
    date_birth  TIMESTAMPTZ,
    uuid        UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

CREATE TABLE IF NOT EXISTS organizations (
    id          BIGSERIAL PRIMARY KEY,
    org_type    SMALLINT NOT NULL DEFAULT 0,
    name        TEXT NOT NULL,
    diocese     TEXT NOT NULL,
    commune     TEXT NOT NULL,
    address     TEXT NOT NULL,
    admin_id    BIGINT NOT NULL REFERENCES users(id),
    uuid        UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT chk_org_type CHECK (org_type BETWEEN 0 AND 255)
    );

CREATE INDEX IF NOT EXISTS idx_organizations_deleted_at ON organizations (deleted_at);

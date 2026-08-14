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

CREATE TABLE IF NOT EXISTS books (
    id            BIGSERIAL PRIMARY KEY,
    org_id        BIGINT NOT NULL REFERENCES organizations(id),
    book_nr       SMALLINT NOT NULL DEFAULT 0,
    date_initial  TIMESTAMPTZ,
    date_final    TIMESTAMPTZ,
    uuid          UUID NOT NULL DEFAULT gen_random_uuid(),
    book_type     SMALLINT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT chk_book_nr   CHECK (book_nr BETWEEN 0 AND 255),
    CONSTRAINT chk_book_type CHECK (book_type BETWEEN 0 AND 255)
    );

CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON books (deleted_at);

CREATE TABLE IF NOT EXISTS register_baptisms (
    id               BIGSERIAL PRIMARY KEY,
    book_id          BIGINT NOT NULL REFERENCES books(id),
    page_number      INTEGER NOT NULL,
    reg_number       INTEGER NOT NULL,
    index_id         BIGINT NOT NULL,
    org_baptism      BIGINT NOT NULL REFERENCES organizations(id),
    baptizer         TEXT NOT NULL,
    date_baptism     TIMESTAMPTZ,
    baptized_name_f  TEXT NOT NULL,
    baptized_name_s  TEXT NOT NULL,
    rut              TEXT NOT NULL,
    date_birth       TIMESTAMPTZ,
    place_birth      TEXT NOT NULL,
    father_name      TEXT NOT NULL,
    father_surname   TEXT NOT NULL,
    mother_name      TEXT NOT NULL,
    mother_surname   TEXT NOT NULL,
    godfather        TEXT NOT NULL,
    godmother        TEXT NOT NULL,
    uuid             UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ
    );

-- index_id intentionally has no foreign key: AddRegisterBaptism (notariat/routes.go)
-- writes it before the matching index_baptisms row exists — the index row is
-- only created afterwards, in updateOrInsertIndexBaptism (notariat/db.go). A hard FK
-- here would reject every register-baptism insert.
CREATE INDEX IF NOT EXISTS idx_register_baptisms_deleted_at ON register_baptisms (deleted_at);

CREATE TABLE IF NOT EXISTS index_baptisms (
    id                BIGSERIAL PRIMARY KEY,
    org_id            BIGINT NOT NULL REFERENCES organizations(id),
    book_id           BIGINT NOT NULL REFERENCES books(id),
    reg_id            BIGINT NOT NULL REFERENCES register_baptisms(id),
    user_surname_f    TEXT NOT NULL,
    user_surname_m    TEXT NOT NULL,
    user_name_first   TEXT NOT NULL,
    user_name_second  TEXT NOT NULL,
    page_number       INTEGER NOT NULL,
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_index_baptisms_deleted_at ON index_baptisms (deleted_at);

CREATE TABLE IF NOT EXISTS certificate_baptisms (
    id               BIGSERIAL PRIMARY KEY,
    org_emisor       BIGINT NOT NULL REFERENCES organizations(id),
    reg_id           BIGINT NOT NULL REFERENCES register_baptisms(id),
    user_validator   BIGINT NOT NULL REFERENCES users(id),
    date_emission    TIMESTAMPTZ,
    date_expiration  TIMESTAMPTZ,
    uuid             UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_certificate_baptisms_deleted_at ON certificate_baptisms (deleted_at);

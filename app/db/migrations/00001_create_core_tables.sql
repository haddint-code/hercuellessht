-- +goose Up
-- SUKAMOTO core schema: users, sessions, packages, package_features, templates.

CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE,
    phone         TEXT        NOT NULL DEFAULT '',
    password_hash TEXT        NOT NULL,
    role          TEXT        NOT NULL CHECK (role IN ('customer','admin')),
    status        TEXT        NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_role ON users (role);

CREATE TABLE sessions (
    id         TEXT PRIMARY KEY,                  -- SCS session token (SHA-256 hashed)
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    data       BYTEA       NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_expires ON sessions (expires_at);

CREATE TABLE packages (
    id               TEXT PRIMARY KEY,            -- 'silver' | 'gold' | 'platinum'
    name             TEXT        NOT NULL,
    tagline          TEXT        NOT NULL DEFAULT '',
    description      TEXT        NOT NULL DEFAULT '',
    price_cents      BIGINT      NOT NULL CHECK (price_cents >= 0),
    currency         TEXT        NOT NULL DEFAULT 'IDR',
    duration_days    INT         NOT NULL DEFAULT 90 CHECK (duration_days > 0),
    max_photos       INT         NOT NULL DEFAULT 0,
    max_videos       INT         NOT NULL DEFAULT 0,
    max_events       INT         NOT NULL DEFAULT 1,
    coming_soon      BOOLEAN     NOT NULL DEFAULT FALSE,
    available        BOOLEAN     NOT NULL DEFAULT TRUE,
    sort_order       INT         NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Feature entitlement flags per package (centralized config; add columns for new features).
CREATE TABLE package_features (
    package_id              TEXT   NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    rsvp_enabled            BOOLEAN NOT NULL DEFAULT FALSE,
    guestbook_enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    digital_envelope        BOOLEAN NOT NULL DEFAULT FALSE,
    whatsapp_button         BOOLEAN NOT NULL DEFAULT FALSE,
    background_music        BOOLEAN NOT NULL DEFAULT FALSE,
    countdown               BOOLEAN NOT NULL DEFAULT FALSE,
    guest_greeting          BOOLEAN NOT NULL DEFAULT FALSE,
    couple_story            BOOLEAN NOT NULL DEFAULT FALSE,
    custom_colors           BOOLEAN NOT NULL DEFAULT FALSE,
    custom_fonts            BOOLEAN NOT NULL DEFAULT FALSE,
    custom_layout           BOOLEAN NOT NULL DEFAULT FALSE,
    custom_cover            BOOLEAN NOT NULL DEFAULT FALSE,
    custom_background       BOOLEAN NOT NULL DEFAULT FALSE,
    custom_sections         BOOLEAN NOT NULL DEFAULT FALSE,
    custom_wording          BOOLEAN NOT NULL DEFAULT FALSE,
    custom_music            BOOLEAN NOT NULL DEFAULT FALSE,
    custom_url              BOOLEAN NOT NULL DEFAULT FALSE,
    multiple_events         BOOLEAN NOT NULL DEFAULT FALSE,
    map_preview             BOOLEAN NOT NULL DEFAULT FALSE,
    event_customization     BOOLEAN NOT NULL DEFAULT FALSE,
    live_rsvp               BOOLEAN NOT NULL DEFAULT FALSE,
    guest_list              BOOLEAN NOT NULL DEFAULT FALSE,
    qr_code                 BOOLEAN NOT NULL DEFAULT FALSE,
    social_sharing_preview  BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (package_id)
);

CREATE TABLE templates (
    id           TEXT PRIMARY KEY,                -- e.g. 'S01', 'G01'
    name         TEXT        NOT NULL,
    category     TEXT        NOT NULL DEFAULT 'general',
    description  TEXT        NOT NULL DEFAULT '',
    thumbnail    TEXT        NOT NULL DEFAULT '',
    package_id   TEXT        NOT NULL REFERENCES packages(id),
    is_active    BOOLEAN     NOT NULL DEFAULT TRUE,
    -- Layout/color/font configuration stored as structured JSONB (per-template defaults).
    config       JSONB       NOT NULL DEFAULT '{}'::jsonb,
    sort_order   INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_templates_package ON templates (package_id);
CREATE INDEX idx_templates_active ON templates (is_active) WHERE is_active;

-- +goose Down
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS package_features;
DROP TABLE IF EXISTS packages;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

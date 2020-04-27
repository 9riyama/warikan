-- +migrate Up

CREATE TABLE users (
  id               SERIAL      PRIMARY KEY
, user_name        TEXT        NOT NULL
, partner_name     TEXT        NOT NULL DEFAULT 'パートナー'
, email            TEXT        NOT NULL
, password         TEXT        NOT NULL
, user_image       TEXT        NOT NULL
, partner_image    TEXT        NOT NULL
, proportion       SMALLINT    NOT NULL DEFAULT 50 --5:5で割り勘
, created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
, updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
, deleted_at       TIMESTAMPTZ
);

CREATE TABLE categories (
  id              SERIAL         PRIMARY KEY
, name            TEXT           NOT NULL
, created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
, updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE payers (
  id              SERIAL         PRIMARY KEY
, name            TEXT           NOT NULL --あなた or パートナー
, created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
, updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE fixed_costs (
  id              SERIAL        PRIMARY KEY
, user_id         SERIAL        NOT NULL REFERENCES users(id)
, category_id     SERIAL        NOT NULL REFERENCES categories(id)
, payer_id        SERIAL        NOT NULL REFERENCES payers(id)
, description     TEXT
, payment_date    TIMESTAMPTZ   NOT NULL
, payment         SERIAL        NOT NULL
, created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
, updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX fixed_costs_user_id_idx        ON fixed_costs (user_id);
CREATE INDEX fixed_costs_category_id_idx    ON fixed_costs (category_id);
CREATE INDEX fixed_costs_payer_id_idx       ON fixed_costs (payer_id);

CREATE TABLE payments (
  id              SERIAL        PRIMARY KEY
, user_id         SERIAL        NOT NULL REFERENCES users(id)
, category_id     SERIAL        NOT NULL REFERENCES categories(id)
, payer_id        SERIAL        NOT NULL REFERENCES payers(id)
, description     TEXT
, payment_date    TIMESTAMPTZ   NOT NULL
, payment         SERIAL        NOT NULL
, created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
, updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX payments_user_id_idx        ON payments (user_id);
CREATE INDEX payments_category_id_idx    ON payments (category_id);
CREATE INDEX payments_payer_id_idx       ON payments (payer_id);

-- +migrate Down

DROP TABLE users;
DROP TABLE categories;
DROP TABLE payers;
DROP TABLE fixed_costs;
DROP TABLE payments;

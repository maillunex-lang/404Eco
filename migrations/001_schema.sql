-- Bashkortostan Economic Simulator schema

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS cities (
    id          TEXT PRIMARY KEY,
    name_ru     TEXT NOT NULL,
    name_ba     TEXT NOT NULL,
    x           REAL NOT NULL,
    y           REAL NOT NULL,
    population  INT  NOT NULL,
    wealth      REAL NOT NULL DEFAULT 1.0,
    industry_bias JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS products (
    id          TEXT PRIMARY KEY,
    name_ru     TEXT NOT NULL,
    category    TEXT NOT NULL,
    base_price  REAL NOT NULL,
    unit        TEXT NOT NULL DEFAULT 'шт'
);

CREATE TABLE IF NOT EXISTS players (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nakama_id   TEXT UNIQUE,
    display_name TEXT NOT NULL,
    cash        DOUBLE PRECISION NOT NULL DEFAULT 5000000,
    reputation  REAL NOT NULL DEFAULT 50,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS firms (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id   UUID REFERENCES players(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    city_id     TEXT NOT NULL REFERENCES cities(id),
    industry    TEXT NOT NULL,
    level       INT NOT NULL DEFAULT 1,
    employees   INT NOT NULL DEFAULT 5,
    capacity    INT NOT NULL DEFAULT 100,
    inventory   JSONB NOT NULL DEFAULT '{}',
    pricing     JSONB NOT NULL DEFAULT '{}',
    quality     REAL NOT NULL DEFAULT 50,
    brand       REAL NOT NULL DEFAULT 10,
    cashflow_d  DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_public   BOOLEAN NOT NULL DEFAULT FALSE,
    share_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    shares_out  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS properties (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id     TEXT NOT NULL REFERENCES cities(id),
    kind        TEXT NOT NULL, -- retail, factory, farm, oil, office, warehouse
    name        TEXT NOT NULL,
    size        INT NOT NULL DEFAULT 1,
    base_rent   DOUBLE PRECISION NOT NULL,
    buy_price   DOUBLE PRECISION NOT NULL,
    owner_firm  UUID REFERENCES firms(id) ON DELETE SET NULL,
    condition   REAL NOT NULL DEFAULT 80
);

CREATE TABLE IF NOT EXISTS market_orders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firm_id     UUID NOT NULL REFERENCES firms(id) ON DELETE CASCADE,
    product_id  TEXT NOT NULL REFERENCES products(id),
    side        TEXT NOT NULL CHECK (side IN ('buy','sell')),
    qty         INT NOT NULL,
    price       DOUBLE PRECISION NOT NULL,
    filled      INT NOT NULL DEFAULT 0,
    status      TEXT NOT NULL DEFAULT 'open',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS stock_holdings (
    player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    firm_id     UUID NOT NULL REFERENCES firms(id) ON DELETE CASCADE,
    shares      INT NOT NULL DEFAULT 0,
    PRIMARY KEY (player_id, firm_id)
);

CREATE TABLE IF NOT EXISTS world_state (
    id          INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    tick        BIGINT NOT NULL DEFAULT 0,
    day         INT NOT NULL DEFAULT 1,
    month       INT NOT NULL DEFAULT 1,
    year        INT NOT NULL DEFAULT 2026,
    oil_price   DOUBLE PRECISION NOT NULL DEFAULT 75,
    grain_index DOUBLE PRECISION NOT NULL DEFAULT 100,
    labor_wage  DOUBLE PRECISION NOT NULL DEFAULT 45000,
    inflation   REAL NOT NULL DEFAULT 0.02,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS news (
    id          BIGSERIAL PRIMARY KEY,
    tick        BIGINT NOT NULL,
    severity    TEXT NOT NULL DEFAULT 'info',
    title       TEXT NOT NULL,
    body        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transactions (
    id          BIGSERIAL PRIMARY KEY,
    tick        BIGINT NOT NULL,
    firm_id     UUID REFERENCES firms(id) ON DELETE SET NULL,
    kind        TEXT NOT NULL,
    amount      DOUBLE PRECISION NOT NULL,
    meta        JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_firms_player ON firms(player_id);
CREATE INDEX IF NOT EXISTS idx_firms_city ON firms(city_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON market_orders(status);
CREATE INDEX IF NOT EXISTS idx_news_tick ON news(tick DESC);

# 404Eco — Капитал Башкортостана

Онлайн экономический симулятор в духе *Capitalism Lab* с действием в **Республике Башкортостан**.

Стек: **React + PixiJS** (клиент), **Go** (игровой тик и экономика), **Nakama** (мультиплеер/RPC), **PostgreSQL** (персистентность и seed).

## Что умеет MVP

- Карта 12 городов РБ (Уфа, Стерлитамак, Салават, Нефтекамск, …) на PixiJS
- Покупка недвижимости и запуск фирм: нефть, АПК, мёд, химия, ритейл, добыча, туризм…
- Производственные цепочки (рецепты входов/выходов), зарплаты, реклама, апгрейды
- Локальный потребительский рынок + ордера buy/sell
- НПЦ-конкуренты, макро (нефть Urals, зерно, новости региона)
- SSE-стрим снапшотов мира; Nakama RPC/match как мультиплеер-мост

Стартовый капитал игрока: **5 000 000 ₽**.

## Быстрый старт (без Docker)

Нужны Go 1.22+ и Node 20+.

```bash
# API (in-memory симуляция — PostgreSQL не обязателен для демо)
cd backend && go mod tidy && USE_MEMORY=true go run ./cmd/server

# Клиент (другой терминал)
cd web && npm install && npm run dev
```

Откройте http://127.0.0.1:5173 — Vite проксирует `/api` на `:8080`.

Или одной командой: `bash scripts/dev.sh`.

## Полный стек (Docker Compose)

```bash
docker compose up --build
```

Сервисы:

| Сервис    | Порт        | Назначение                          |
|-----------|-------------|-------------------------------------|
| web       | 5173→80     | React UI                            |
| game-api  | 8080        | Экономический движок                |
| nakama    | 7349/7350/7351 | Auth, RPC, lobby match           |
| postgres  | 5432        | Игровая БД (`migrations/*.sql`)     |
| nakama-db | internal    | Служебная БД Nakama                 |

Консоль Nakama: http://127.0.0.1:7351 (`admin` / `adminpass`).

## Архитектура

```
web (React/Pixi)
   │  REST + SSE                │  Nakama JS client (опционально)
   ▼                            ▼
game-api (Go tick engine) ←—— nakama/modules/bashkir_eco.js (RPC)
   │
   ▼
PostgreSQL (schema + seed городов/товаров/площадок)
```

- `backend/internal/economy` — тик: производство, матчинг ордеров, ритейл, НПЦ, макрособытия
- `nakama/modules` — `eco_register`, `eco_create_firm`, `eco_firm_action`, `eco_place_order`, match `bashkir_lobby`
- `migrations/` — схема и сиды РБ (города, продукты, объекты)

## Геймплей (кратко)

1. Войдите под именем предпринимателя.
2. Клик по городу на карте → купите свободный объект → выберите отрасль.
3. Нанимайте персонал, крутите рекламу, выставляйте ордера на бирже.
4. Следите за нефтью, зерном и лентой новостей республики.

## Переменные окружения

**game-api:** `DATABASE_URL`, `HTTP_ADDR`, `TICK_INTERVAL_MS`, `CORS_ORIGIN`, `NAKAMA_URL`, `NAKAMA_HTTP_KEY`

**web:** `VITE_API_URL`, `VITE_NAKAMA_HOST`, `VITE_NAKAMA_PORT`, `VITE_NAKAMA_KEY`

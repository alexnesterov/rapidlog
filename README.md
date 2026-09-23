# RapidLog

Цифровой bullet journal для быстрых заметок, задач и событий. Позволяет вести
журнал дня — задачи можно завершать, переносить на другой день или отменять,
как в бумажном bullet journal.

Backend на Go + Postgres, фронтенд на React встраивается в бинарник `gate`.

Состоит из двух сервисов в одном Go-модуле:

- **gate** — единственная точка входа для браузера: REST API, раздача SPA и весь
  домен bullets.
- **identity** — резолюция анонимной сессии в пользователя; общается с `gate`
  по gRPC и владеет собственной БД.

## Запуск

Полный стек в Docker (gate + identity + Postgres + Adminer):

```sh
cp .env.example .env
docker compose up
```

Приложение доступно на `http://localhost:8080`, Adminer — на `:8081`,
gRPC `identity` — на `:50051`, Postgres — на `:5432`. Вторая база `identity`
создаётся в том же Postgres при первом старте на пустом volume.

Аутентификации нет: анонимная привязка к пользователю происходит через
cookie `session_id`, которую `gate` выставляет сам при первом запросе (сессия
резолвится вызовом к `identity`). Новому пользователю при первом просмотре
списка создаются три демонстрационных bullet'а.

### Без Docker для самих сервисов

Сервисам нужен Postgres и собранный фронтенд — `gate` отдаёт его из `web/dist`,
поэтому без сборки `go run`/`go build` для него падают. Проще всего поднять в
Docker только БД, остальное запустить локально:

```sh
cp .env.example .env
docker compose up -d db
cd web && npm ci && npm run build && cd ..
go run ./cmd/identity     # gRPC на :50051
go run ./cmd/gate         # HTTP на :8080
```

`identity` нужно запустить до `gate`: сессия резолвится gRPC-вызовом на
каждый запрос (кроме `/health`), без `identity` они отвечают 500. Если одновременно
работает docker-контейнер `identity`, он занимает порт 50051 — остановите его
(`docker compose stop identity`).

Миграции применяются автоматически при старте каждого сервиса.

### Конфигурация

Только переменные окружения, файлов конфигурации нет. Дефолты рассчитаны на
запуск на хосте против портов, проброшенных из compose (`localhost:5432`,
`localhost:50051`); в compose их переопределяет `.env` (см. `.env.example`).

| Сервис   | Префикс      | Переменные                                                                  |
| -------- | ------------ | --------------------------------------------------------------------------- |
| gate     | `GATE_`      | `GATE_DSN`, `GATE_PORT`, `GATE_IDENTITY`                                    |
| identity | `IDENTITY_`  | `IDENTITY_DB_DSN`, `IDENTITY_GRPC_PORT`                                     |

## Разработка

```sh
go build -v ./...                      # сборка
go test -v ./...                       # тесты
go test -race -v ./...                 # тесты с race detector (гоняются в CI)
go test -run TestName ./internal/...   # один тест
golangci-lint run                      # линтер (в CI не запускается)
```

Frontend (`web/`):

```sh
npm run dev       # dev-сервер vite (проксирует /api на localhost:8080)
npm run build     # сборка в web/dist
npm run lint      # oxlint
```

Моки генерируются [mockery](https://vektra.github.io/mockery/) по `.mockery.yml`:
после изменения интерфейса в `port`- или `httpapi`-пакетах сервиса
перегенерировать командой `mockery`.

gRPC-контракт `identity` (`proto/identity/v1/identity.proto`) описан через
[buf](https://buf.build/); после правки `.proto` — `buf generate`
(сгенерированный код в `gen/` коммитится).

Юнит-тестами покрыты usecase и хендлеры; Postgres-репозитории, загрузка
конфигурации и gRPC-клиент проверяются вручную через `docker compose`.

## API

Формат ответа единый: `{"data": ...}` при успехе, `{"data": null, "error":
{"code", "message"}}` при ошибке. Подробные контракты — в
[`docs/api-design.md`](docs/api-design.md).

```
GET    /health
POST   /api/bullets
GET    /api/bullets
POST   /api/bullets/{id}/complete
POST   /api/bullets/{id}/migrate
POST   /api/bullets/{id}/cancel
```

## Структура

```text
cmd/gate                                     — точка входа gate
cmd/identity                                 — точка входа identity
internal/service/gate                        — фасад gate: Run(ctx, logger)
internal/service/gate/internal/adapter       — HTTP-хендлеры и middleware (Recovery, Logging, Session)
internal/service/gate/internal/config        — конфигурация из env (Viper)
internal/service/gate/internal/domain        — сущности, порты, use cases доменного слоя bullets
internal/service/gate/internal/infrastructure — репозитории на pgx, транзакции, gRPC-клиент identity
internal/service/gate/migrations             — SQL-миграции gate (golang-migrate)
internal/service/identity                    — фасад identity: Run(ctx, logger)
internal/service/identity/internal           — gRPC-хендлер, домен пользователя, репозиторий, конфиг
internal/service/identity/migrations         — SQL-миграции identity
proto, gen                                   — gRPC-контракт identity и сгенерированный код
docs/api-design.md                           — контракты API, форматы запросов/ответов
web                                          — React-интерфейс (Vite + TypeScript)
```

Реализация каждого сервиса лежит в `internal/service/<name>/internal/...`,
поэтому импортировать её может только код внутри `internal/service/<name>/` —
сервисы не могут обратиться к внутренностям друг друга.

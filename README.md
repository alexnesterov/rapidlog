# RapidLog

Цифровой bullet journal для быстрых заметок, задач и событий. Позволяет вести
журнал дня — задачи можно завершать, переносить на другой день или отменять,
как в бумажном bullet journal.

Backend на Go + Postgres, фронтенд на React встраивается в бинарник `gateway`.

Состоит из трёх сервисов в одном Go-модуле:

- **gateway** — единственная точка входа для браузера: REST API, раздача SPA.
  Через gRPC общается с `identity` и `journal`.
- **identity** — резолюция анонимной сессии в пользователя; владеет собственной БД.
- **journal** — домен bullets, хранит и обрабатывает записи; владеет собственной БД.

## Запуск

Полный стек в Docker (gateway + identity + journal + Postgres + Adminer):

```sh
cp .env.example .env
docker compose up
```

Приложение доступно на `http://localhost:8080`, Adminer — на `:8081`,
gRPC `identity` — на `:50051`, gRPC `journal` — на `:50052`, Postgres — на `:5432`.
Дополнительные базы `identity` и `journal` создаются в том же Postgres при первом
старте на пустом volume.

Аутентификации нет: анонимная привязка к пользователю происходит через
cookie `session_id`, которую `gateway` выставляет сам при первом запросе (сессия
резолвится вызовом к `identity`). Новому пользователю при первом просмотре
списка создаются три демонстрационных bullet'а в `journal`.

### Без Docker для самих сервисов

Сервисам нужен Postgres и собранный фронтенд — `gateway` отдаёт его из `web/dist`,
поэтому без сборки `go run`/`go build` для него падают. Проще всего поднять в
Docker только БД, остальное запустить локально:

```sh
cp .env.example .env
docker compose up -d db
cd web && npm ci && npm run build && cd ..
go run ./cmd/identity     # gRPC на :50051
go run ./cmd/journal      # gRPC на :50052
go run ./cmd/gateway      # HTTP на :8080
```

`identity` и `journal` нужно запустить до `gateway`: сессия резолвится gRPC-вызовом
на каждый запрос (кроме `/health`), без `identity` и `journal` они отвечают 500.
Если одновременно работают docker-контейнеры `identity` или `journal`, они занимают
порты 50051/50052 — остановите их (`docker compose stop identity journal`).

Миграции применяются автоматически при старте каждого сервиса.

### Конфигурация

Только переменные окружения, файлов конфигурации нет. Дефолты рассчитаны на
запуск на хосте против портов, проброшенных из compose (`localhost:5432`,
`localhost:50051`); в compose их переопределяет `.env` (см. `.env.example`).

| Сервис   | Префикс      | Переменные                                                                       |
| -------- | ------------ | -------------------------------------------------------------------------------- |
| gateway  | `GATEWAY_`   | `GATEWAY_PORT`, `GATEWAY_IDENTITY`, `GATEWAY_JOURNAL`                            |
| identity | `IDENTITY_`  | `IDENTITY_DSN`, `IDENTITY_PORT`                                                  |
| journal  | `JOURNAL_`    | `JOURNAL_DSN`, `JOURNAL_PORT`                                                    |

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

gRPC-контракты (`proto/identity/v1/identity.proto`, `proto/journal/v1/journal.proto`)
описаны через [buf](https://buf.build/); после правки `.proto` — `buf generate`
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
cmd/gateway                                  — точка входа gateway
cmd/identity                                 — точка входа identity
cmd/journal                                  — точка входа journal
internal/service/gateway                     — фасад gateway: Run(ctx, logger)
internal/service/gateway/internal            — HTTP-хендлеры, middleware, конфигурация
internal/service/identity                    — фасад identity: Run(ctx, logger)
internal/service/identity/internal           — gRPC-хендлер, домен пользователя, репозиторий, конфиг
internal/service/identity/migrations         — SQL-миграции identity
internal/service/journal                     — фасад journal: Run(ctx, logger)
internal/service/journal/internal            — gRPC-хендлер, домен bullets, репозиторий, конфиг
internal/service/journal/migrations           — SQL-миграции journal
proto, gen                                   — gRPC-контракты (identity, journal) и сгенерированный код
docs/api-design.md                           — контракты API, форматы запросов/ответов
web                                          — React-интерфейс (Vite + TypeScript)
```

Реализация каждого сервиса лежит в `internal/service/<name>/internal/...`,
поэтому импортировать её может только код внутри `internal/service/<name>/` —
сервисы не могут обратиться к внутренностям друг друга.

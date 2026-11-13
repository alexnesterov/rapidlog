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
порты 50051/50052 — остановите их (`docker compose stop identity journal`) или
задайте другие `IDENTITY_PORT`/`JOURNAL_PORT` и направьте на них `gateway` через
`GATEWAY_IDENTITY`/`GATEWAY_JOURNAL`.

Миграции применяются автоматически при старте `identity` и `journal`.

### Продакшн

`docker-compose.yaml` (не путать с локальным `compose.yaml`) поднимает только
`gateway`, `identity` и `journal` — без Postgres и Adminer. Значения переменных
подставляются из окружения хоста (`${VAR}`), DSN указывают на внешний Postgres.

### Конфигурация

Только переменные окружения, файлов конфигурации нет. Дефолты рассчитаны на
запуск на хосте против портов, проброшенных из compose (`localhost:5432`,
`localhost:50051`, `localhost:50052`); в compose их переопределяет `.env`
(см. `.env.example`).

| Сервис   | Переменная                       | Дефолт                                              |
| -------- | -------------------------------- | --------------------------------------------------- |
| gateway  | `GATEWAY_PORT`                   | `8080`                                              |
|          | `GATEWAY_IDENTITY`               | `localhost:50051`                                   |
|          | `GATEWAY_JOURNAL`                | `localhost:50052`                                   |
|          | `GATEWAY_HTTP_READ_TIMEOUT`      | `5s`                                                |
|          | `GATEWAY_HTTP_WRITE_TIMEOUT`     | `10s`                                               |
|          | `GATEWAY_HTTP_IDLE_TIMEOUT`      | `15s`                                               |
|          | `GATEWAY_SESSION_COOKIE_NAME`    | `session_id`                                        |
|          | `GATEWAY_SESSION_COOKIE_TTL`     | `8760h` (год)                                       |
|          | `GATEWAY_SESSION_COOKIE_SECURE`  | `true`                                              |
| identity | `IDENTITY_DSN`                   | `postgres://rapidlog:rapidlog@localhost:5432/identity` |
|          | `IDENTITY_PORT`                  | `50051`                                             |
| journal  | `JOURNAL_DSN`                    | `postgres://rapidlog:rapidlog@localhost:5432/journal`  |
|          | `JOURNAL_PORT`                   | `50052`                                             |

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

Те же операции обёрнуты в [Taskfile](https://taskfile.dev/) (`Taskfile.yml`, список —
`task --list`): `task build`, `task test`, `task test-cov` (сквозное покрытие),
`task test-tdd DIR=./internal/...`, `task mock`, `task buf`, `task web`,
`task up`/`task down` (`docker compose up -d --build` / `down`).

Моки генерируются [mockery](https://vektra.github.io/mockery/) v3 по `.mockery.yml`
(`internal/pkg/...`, `port`-пакеты `identity` и `journal`, `httpapi` в `gateway`;
моки лежат в соседнем `mocks/`): после изменения интерфейса перегенерировать
командой `mockery`. Устаревшие `mock_*.go` после переименования интерфейса
mockery не удаляет — их нужно удалить вручную.

gRPC-контракты (`api/proto/identity/v1/identity.proto`,
`api/proto/journal/v1/journal.proto`) описаны через [buf](https://buf.build/);
после правки `.proto` — `task buf` (ставит buf и плагины `protoc-gen-go`/`protoc-gen-go-grpc`
в `bin/` и запускает `buf generate`; сгенерированный код в `api/gen/`
коммитится, в CI `buf` не запускается).

Юнит-тестами покрыты usecase и gRPC-хендлеры `identity` и `journal`.
Postgres-репозитории, загрузка конфигурации и gRPC-клиенты проверяются вручную
через `docker compose` + `curl`/`grpcurl`. REST-хендлеры `gateway` пока без тестов.

Для ручных запросов есть [Bruno](https://www.usebruno.com/)-коллекция в
`api/bruno/`: `bullets/` — REST через `gateway`, `journal/` — gRPC напрямую к
`journal`; окружения `Local` и `Docker`.

## API

Формат ответа единый: `{"data": ...}` при успехе, `{"data": null, "error":
{"code", "message"}}` при ошибке. Тела запросов и ответов — protojson-сообщения
`journal` как есть (camelCase, enum'ы вида `SIGNIFIER_OPEN`); ошибки `journal`
пока не маппятся в HTTP-коды и возвращаются как 500. Подробные контракты — в
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
internal/pkg                                 — общий код: transactor (pgx-транзакции), protoenum
internal/service/gateway                     — фасад gateway: Run(ctx, logger)
internal/service/gateway/config              — конфигурация gateway
internal/service/gateway/pkg/respond         — конверт ответа {"data", "error"}
internal/service/gateway/internal            — HTTP-хендлеры, middleware
internal/service/identity                    — фасад identity: Run(ctx, logger)
internal/service/identity/config             — конфигурация identity
internal/service/identity/pkg/identitygrpc   — gRPC-клиент identity (используется gateway)
internal/service/identity/internal           — gRPC-хендлер, домен пользователя, репозиторий
internal/service/identity/migrations         — SQL-миграции identity
internal/service/journal                     — фасад journal: Run(ctx, logger)
internal/service/journal/config              — конфигурация journal
internal/service/journal/pkg/journalgrpc     — gRPC-клиент journal (используется gateway)
internal/service/journal/internal            — gRPC-хендлеры, домен bullets, репозиторий
internal/service/journal/migrations          — SQL-миграции journal
api/proto, api/gen                           — gRPC-контракты (identity, journal) и сгенерированный код
api/bruno                                    — Bruno-коллекция для ручных запросов
docs/api-design.md                           — контракты API, форматы запросов/ответов
web                                          — React-интерфейс (Vite + TypeScript)
compose.yaml, docker-compose.yaml            — локальный полный стек / продакшн
```

Реализация каждого сервиса лежит в `internal/service/<name>/internal/...`,
поэтому импортировать её может только код внутри `internal/service/<name>/` —
сервисы не могут обратиться к внутренностям друг друга. Наружу сервис открывает
только `config/` и `pkg/`: `gateway` ходит в `identity` и `journal` через их
gRPC-клиенты из `pkg/`.

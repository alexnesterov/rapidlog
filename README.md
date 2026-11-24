# RapidLog

Цифровой bullet journal для быстрых заметок, задач и событий. Позволяет вести
журнал дня — задачи можно завершать, переносить на другой день или отменять,
как в бумажном bullet journal.

Backend на Go + Postgres, фронтенд на React встраивается в бинарник `gateway`.

Состоит из двух сервисов в одном Go-модуле и [Ory Kratos](https://www.ory.sh/kratos/)
для аккаунтов:

- **gateway** — единственная точка входа для браузера: REST API, раздача SPA,
  проверка сессии Kratos и прокси его публичного API на `/auth/*`. Через gRPC
  общается с `journal`.
- **journal** — домен bullets, хранит и обрабатывает записи; владеет собственной БД.
- **kratos** — регистрация и вход по email + паролю, сессии; своя БД `kratos`.
  Конфигурация — `deploy/kratos/`.

## Запуск

Полный стек в Docker (gateway + journal + kratos + Postgres):

```sh
cp .env.example .env
docker compose up
```

`docker compose` сам подхватывает `compose.override.yaml` поверх `compose.yaml`:
он подключает Postgres из `infra.yaml`, пробрасывает порты `journal` и Kratos на
хост и запускает Kratos с `--dev` (cookie без Secure по `http://localhost`).

Приложение доступно на `http://localhost:8080`,
gRPC `journal` — на `:50052`, публичный API Kratos — на `:4433`, Postgres — на `:5432`.
Дополнительные базы `journal` и `kratos` создаются в том же Postgres при первом
старте на пустом volume; схему Kratos накатывает одноразовый контейнер
`kratos-migrate` перед стартом `kratos`.

Вход — по email и паролю через Kratos. Браузер работает с Kratos только через
`gateway` (`/auth/*`), поэтому cookie сессии живёт на том же origin, что и SPA.
`gateway` на каждый запрос к `/api/*` проверяет сессию вызовом Kratos
`/sessions/whoami`: без сессии — `401`, с сессией ID identity Kratos становится
`user_id` в `journal`. Новому пользователю при первом просмотре списка создаются
три демонстрационных bullet'а.

### Без Docker для самих сервисов

Сервисам нужен Postgres, Kratos и собранный фронтенд — `gateway` отдаёт его из
`web/dist`, поэтому без сборки `go run`/`go build` для него падают. Проще всего
поднять в Docker БД и Kratos, остальное запустить локально:

```sh
cp .env.example .env
docker compose up -d postgres kratos
cd web && npm ci && npm run build && cd ..
go run ./cmd/journal      # gRPC на :50052
go run ./cmd/gateway      # HTTP на :8080
```

`journal` нужно запустить до `gateway`, иначе операции с bullet'ами отвечают 500;
без Kratos любой запрос к `/api/*` отвечает 500 (сессию проверить нечем). Если
одновременно работает docker-контейнер `journal`, он занимает порт 50052 —
остановите его (`docker compose stop journal`) или задайте другой `JOURNAL_PORT`
и направьте на него `gateway` через `GATEWAY_JOURNAL`.

Миграции `journal` применяются автоматически при его старте.

### Продакшн

Продакшн — это `compose.yaml` без override: только `gateway`, `journal`,
`kratos-migrate` и `kratos`, без Postgres. Значения переменных подставляются из
окружения хоста или серверного `.env` (`${VAR}`), DSN указывают на внешний Postgres.

Kratos в продакшене запускается без `--dev`, поэтому публичный URL приложения
(`KRATOS_APP_URL`) должен быть на HTTPS — иначе браузер не сохранит Secure-cookie
сессии.

### Конфигурация

Сервисы на Go — только переменные окружения, файлов конфигурации нет. Дефолты
рассчитаны на запуск на хосте против портов, проброшенных из compose
(`localhost:5432`, `localhost:50052`, `localhost:4433`); в compose их переопределяет
`.env` (см. `.env.example`). Kratos читает `deploy/kratos/kratos.yml`, а URL, DSN и
секреты в `compose.yaml` переопределяются переменными `KRATOS_*`.

| Сервис   | Переменная                       | Дефолт                                              |
| -------- | -------------------------------- | --------------------------------------------------- |
| gateway  | `GATEWAY_PORT`                   | `8080`                                              |
|          | `GATEWAY_JOURNAL`                | `localhost:50052`                                   |
|          | `GATEWAY_KRATOS_PUBLIC_URL`      | `http://localhost:4433`                             |
|          | `GATEWAY_HTTP_READ_TIMEOUT`      | `5s`                                                |
|          | `GATEWAY_HTTP_WRITE_TIMEOUT`     | `10s`                                               |
|          | `GATEWAY_HTTP_IDLE_TIMEOUT`      | `15s`                                               |
| journal  | `JOURNAL_DSN`                    | `postgres://rapidlog:rapidlog@localhost:5432/journal`  |
|          | `JOURNAL_PORT`                   | `50052`                                             |
| kratos   | `KRATOS_DSN`                     | `postgres://rapidlog:rapidlog@postgres:5432/kratos?sslmode=disable` (compose) |
|          | `KRATOS_SECRETS_COOKIE`          | dev-значение из `kratos.yml`, ≥ 16 символов         |
|          | `KRATOS_SECRETS_CIPHER`          | dev-значение из `kratos.yml`, ровно 32 символа      |
|          | `KRATOS_APP_URL`                 | `http://localhost:8080`                             |

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
npm run dev       # dev-сервер vite (проксирует /api и /auth на localhost:8080)
npm run build     # сборка в web/dist
npm run lint      # oxlint
```

Те же операции обёрнуты в [Taskfile](https://taskfile.dev/) (`Taskfile.yml`, список —
`task --list`): `task build`, `task test`, `task test-cov` (сквозное покрытие),
`task test-tdd DIR=./internal/...`, `task mock`, `task buf`, `task web`,
`task up`/`task down` (`docker compose up -d --build` / `down`).

Моки генерируются [mockery](https://vektra.github.io/mockery/) v3 по `.mockery.yml`
(`internal/pkg/...`, `usecase` и `controller/grpcapi` в `journal`, `httpapi` в `gateway`;
моки лежат в соседнем `mocks/`): после изменения интерфейса перегенерировать
командой `task mock` (ставит mockery нужной версии в `bin/` и запускает его). Устаревшие `mock_*.go` после переименования интерфейса
mockery не удаляет — их нужно удалить вручную.

gRPC-контракт `api/proto/journal/v1/journal.proto` описан через [buf](https://buf.build/);
после правки `.proto` — `task buf` (ставит buf и плагины `protoc-gen-go`/`protoc-gen-go-grpc`
в `bin/` и запускает `buf generate`; сгенерированный код в `api/gen/`
коммитится, в CI `buf` не запускается).

Юнит-тестами покрыты usecase и gRPC-хендлеры `journal`, а в `gateway` — middleware
проверки сессии и клиент Kratos (на `httptest`). Postgres-репозитории, загрузка
конфигурации, gRPC-клиенты и SPA-fallback проверяются вручную через `docker compose` +
`curl`/`grpcurl`. REST-хендлеры `gateway` пока без тестов.

Для ручных запросов есть [Bruno](https://www.usebruno.com/)-коллекция в
`api/bruno/`: `auth/` — регистрация и вход через Kratos (сохраняют токен сессии в
`sessionToken`), `bullets/` — REST через `gateway` с заголовком `X-Session-Token`,
`journal/` — gRPC напрямую к `journal`; окружения `Local` и `Docker` (учётные
данные — переменные `email`/`password`).

## API

Формат ответа единый: `{"data": ...}` при успехе, `{"data": null, "error":
{"code", "message"}}` при ошибке. Тела запросов и ответов — protojson-сообщения
`journal` как есть (camelCase, enum'ы вида `SIGNIFIER_OPEN`); ошибки `journal`
пока не маппятся в HTTP-коды и возвращаются как 500. Подробные контракты — в
[`docs/api-design.md`](docs/api-design.md).

Все `/api/*` требуют сессию Kratos — cookie браузера или заголовок
`X-Session-Token` (API-клиенты); без неё — `401`. `/auth/*` — публичный API Kratos
(self-service flows регистрации/входа/выхода, `sessions/whoami`), проксируется как есть.

```
GET    /health
ANY    /auth/*
POST   /api/bullets
GET    /api/bullets
POST   /api/bullets/{id}/complete
POST   /api/bullets/{id}/migrate
POST   /api/bullets/{id}/cancel
```

## Структура

```text
cmd/gateway                                  — точка входа gateway
cmd/journal                                  — точка входа journal
internal/pkg                                 — общий код: transactor (pgx-транзакции), protoenum
internal/service/gateway                     — фасад gateway: Run(ctx, logger)
internal/service/gateway/config              — конфигурация gateway
internal/service/gateway/pkg/respond         — конверт ответа {"data", "error"}
internal/service/gateway/internal            — HTTP-хендлеры, middleware, адаптер Kratos
internal/service/journal                     — фасад journal: Run(ctx, logger)
internal/service/journal/config              — конфигурация journal
internal/service/journal/pkg/journalgrpc     — gRPC-клиент journal (используется gateway)
internal/service/journal/internal            — gRPC-хендлеры, домен bullets, репозиторий
internal/service/journal/migrations          — SQL-миграции journal
api/proto, api/gen                           — gRPC-контракт journal и сгенерированный код
api/bruno                                    — Bruno-коллекция для ручных запросов
deploy/kratos                                — конфиг Ory Kratos и схема identity
docs/api-design.md                           — контракты API, форматы запросов/ответов
web                                          — React-интерфейс (Vite + TypeScript)
compose.yaml                                 — продакшн (без Postgres)
compose.override.yaml, infra.yaml            — локальная надстройка и Postgres для неё
```

Реализация каждого сервиса лежит в `internal/service/<name>/internal/...`,
поэтому импортировать её может только код внутри `internal/service/<name>/` —
сервисы не могут обратиться к внутренностям друг друга. Наружу сервис открывает
только `config/` и `pkg/`: `gateway` ходит в `journal` через gRPC-клиент из его
`pkg/`, а в Kratos — по HTTP через свой адаптер `internal/adapter/kratos`.

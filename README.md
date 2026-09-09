# RapidLog

HTTP API для ведения bullet journal: заметки (bullets) с типом
(задача/событие/заметка) и статусом (открыта/выполнена/перенесена/отменена).
Backend на Go + Postgres, фронтенд на React встраивается в тот же бинарник.

## Запуск

Нужен Postgres (см. `db.dsn` в `configs/config.yaml` или `.env.example`) и
собранный фронтенд — сервер отдаёт его из `web/dist`, поэтому без сборки
`go run`/`go build` падают:

```sh
cd web && npm ci && npm run build
cd .. && go run ./cmd/app
```

Сервер стартует на `:1508`, применяет миграции автоматически при старте и
раздаёт SPA по `/`, API — по `/api/*`.

Либо полный стек в Docker (app + Postgres + Adminer на :8080):

```sh
cp .env.example .env
docker compose up
```

App будет доступен на хостовом порту 1509.

Аутентификации нет: анонимная привязка к пользователю происходит через
cookie `session_id`, которую сервер выставляет сам при первом запросе.

## Разработка

```sh
go build -v ./...                      # сборка
go test -v ./...                       # тесты
go test -race -v ./...                 # тесты с race detector (гоняются в CI)
go test -run TestName ./internal/...   # один тест
```

Frontend (`web/`):

```sh
npm run dev       # dev-сервер vite
npm run build     # сборка в web/dist
npm run lint      # oxlint
```

Моки (`internal/domain/port/mocks`, `internal/adapter/httpapi/mocks`)
генерируются [mockery](https://vektra.github.io/mockery/) по `.mockery.yml`:
после изменения интерфейса в `internal/domain/port` или
`internal/adapter/httpapi` — перегенерировать командой `mockery`.

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
cmd/app                             — точка входа сервера
internal/adapter/httpapi            — HTTP-хендлеры
internal/adapter/httpapi/middleware — Recovery, Logging, Session
internal/adapter/httpapi/mocks      — моки httpapi-интерфейсов, сгенерированные mockery
internal/config                     — конфигурация приложения (Viper)
internal/domain/entity              — доменные сущности и их правила переходов
internal/domain/port                — интерфейсы (порты) доменного слоя
internal/domain/port/mocks          — моки портов, сгенерированные mockery
internal/domain/usecase             — реализации доменных сервисов (use cases)
internal/infrastructure/postgres    — репозитории на pgx, транзакции, миграции
migrations                          — SQL-миграции (golang-migrate), встраиваются в бинарник
docs/api-design.md                  — контракты API, форматы запросов/ответов
web                                 — React-интерфейс (Vite + TypeScript)
```

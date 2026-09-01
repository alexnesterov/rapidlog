# RapidLog

Простой HTTP-сервер на Go.

## Запуск

```sh
go run ./cmd/app
```

Сервер стартует на `:1508`.

## Веб-интерфейс

React-интерфейс находится в `web/` и после сборки встраивается в бинарник
API. Go раздаёт его по `/`, а маршруты `/api/*` продолжают обрабатываться
backend'ом.

Для локального запуска сначала соберите frontend, затем API:

```sh
cd web && npm ci && npm run build
cd .. && go run ./cmd/app
```

Docker-сборка выполняет эти шаги автоматически; отдельный frontend-сервис
или Node.js в итоговом образе не нужны.

## Структура

```text
cmd/app                             — точка входа сервера
cmd/client                          — вспомогательный HTTP-клиент для ручных запросов
internal/adapter/httpapi            — HTTP-хендлеры и middleware
internal/config                     — конфигурация приложения
internal/domain/entity              — доменные сущности
internal/domain/port                — интерфейсы (порты) доменного слоя
internal/domain/port/mocks          — моки портов, сгенерированные mockery
internal/domain/usecase             — реализации доменных сервисов (use cases)
internal/infrastructure/postgres    — реализации репозиториев на pgx, миграции
migrations                          — SQL-миграции (golang-migrate), встраиваются в бинарник
web                                 — React-интерфейс (Vite + TypeScript)
```

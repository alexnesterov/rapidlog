# RapidLog

Простой HTTP-сервер на Go.

## Запуск

```sh
go run ./cmd/api
```

Сервер стартует на `:1508`.

## Веб-интерфейс

React-интерфейс находится в `web/` и после сборки встраивается в бинарник
API. Go раздаёт его по `/`, а маршруты `/api/*` продолжают обрабатываться
backend'ом.

Для локального запуска сначала соберите frontend, затем API:

```sh
cd web && npm ci && npm run build
cd .. && go run ./cmd/api
```

Docker-сборка выполняет эти шаги автоматически; отдельный frontend-сервис
или Node.js в итоговом образе не нужны.

## Структура

```text
cmd/api                            — точка входа
internal/adapter/httpapi           — HTTP-хендлеры
internal/config                    — конфигурация приложения
internal/domain/entity             — доменные сущности
internal/domain/port               — интерфейсы (порты) доменного слоя
internal/domain/port/mocks         — моки портов, сгенерированные mockery
internal/domain/usecase            — реализации доменных сервисов (use cases)
internal/infrastructure/postgresql — реализации репозиториев
```

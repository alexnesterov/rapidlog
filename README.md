# RapidLog

Простой HTTP-сервер на Go.

## Запуск

```sh
go run ./cmd/api
```

Сервер стартует на `:1508`.

## Структура

```text
cmd/api                            — точка входа
internal/adapter/httpapi           — HTTP-хендлеры
internal/config                    — конфигурация приложения
internal/domain/entity             — доменные сущности
internal/domain/port               — интерфейсы (порты) доменного слоя
internal/infrastructure/postgresql — реализации репозиториев
```

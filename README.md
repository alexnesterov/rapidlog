# RapidLog

Простой HTTP-сервер на Go.

## Запуск

```sh
go run ./cmd/app
```

Сервер стартует на `:1508`.

## Структура

```text
cmd/app                            — точка входа
internal/adapter/httpapi           — HTTP-хендлеры
internal/config                    — конфигурация приложения
internal/domain/entity             — доменные сущности
internal/domain/port               — интерфейсы (порты) доменного слоя
internal/infrastructure/repository — реализации репозиториев
```

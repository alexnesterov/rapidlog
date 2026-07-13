# DotLine

Простой HTTP-сервер на Go.

## Запуск

```sh
go run ./cmd/app
```

Сервер стартует на `:8080`.

## Структура

```text
cmd/app                    — точка входа
internal/adapter/httpapi   — HTTP-хендлеры
```

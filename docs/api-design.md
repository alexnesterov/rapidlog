# API — проектирование

> Этот контракт обслуживает сервис `gateway` (`internal/service/gateway/`) — единственная
> точка входа для браузера. Сессию он резолвит gRPC-вызовом к `identity`, операции с bullet'ами
> проксирует gRPC-вызовами к `journal`. Тела ответов — это gRPC-ответы `journal`,
> сериализованные `protojson` (`api/proto/journal/v1/journal.proto`), поэтому поля в
> camelCase, а значения enum'ов — полные имена из proto (`KIND_TASK`, `SIGNIFIER_OPEN`).

## 1. Сущности

### Bullet

| Поле          | Тип       | Описание                                                          |
| ------------- | --------- | ----------------------------------------------------------------- |
| id            | UUID      |                                                                   |
| userId        | UUID      | анонимный владелец (см. привязку к сессии)                       |
| kind          | string    | `KIND_TASK` \| `KIND_EVENT` \| `KIND_NOTE`, обязательное |
| signifier     | string    | `SIGNIFIER_OPEN` \| `SIGNIFIER_COMPLETED` \| `SIGNIFIER_CANCELLED` |
| content       | string    | обязательное, ≤200 символов                                       |
| createdAt     | timestamp | RFC 3339, UTC                                                     |
| updatedAt     | timestamp | RFC 3339, UTC                                                     |

По аналогии с нотацией Bullet Journal: `kind` описывает разновидность
записи (задача, событие, заметка), `signifier` — её состояние (открыта,
выполнена, отменена). Реализованы пока переходы
`open` → `completed` и `open` → `cancelled` (см. `POST
/api/bullets/{id}/complete` и `POST /api/bullets/{id}/cancel`).
Перенос (`POST /api/bullets/{id}/migrate`) статус не меняет: открытая
задача переезжает в текущий день. Ниже в тексте
значения сокращены (`open`, `task`), в JSON они всегда полные.

**Операции:** Create, List, Complete, Migrate, Cancel.

### Collection

> Реализация отложена: сначала MVP на `Bullet`, `Collection`
> подключается позже.

По аналогии с Bullet Journal: каждый bullet живёт в какой-то
коллекции — произвольном именованном списке.

| Поле       | Тип       | Описание                    |
| ---------- | --------- | --------------------------- |
| id         | UUID      |                             |
| topic      | string    | обязательное, ≤100 символов |
| created_at | timestamp |                             |
| updated_at | timestamp |                             |

**Операции:** Create, List, Get, Update, Delete.

---

## 2. Общий формат ответов

**Успех — одиночный объект:**

```json
{
  "data": {
    /* сущность */
  }
}
```

**Успех — список:**

```json
{
  "data": [
    /* Bullet[] */
  ]
}
```

**Ошибка:**

```json
{
  "data": null,
  "error": {
    "code": 400,
    "message": "bullet content is required"
  }
}
```

`error.code` — числовой HTTP-статус (дублирует статус-код ответа),
отдельного строкового кода ошибки нет.

| Код | Когда                                                 |
| --- | ----------------------------------------------------- |
| 200 | успешный GET/PUT/POST без создания ресурса            |
| 201 | ресурс создан                                         |
| 204 | успешное удаление, тело пустое                        |
| 400 | некорректный запрос (невалидный JSON)                 |
| 401 | в контексте запроса нет пользователя                  |
| 404 | ресурс не найден                                      |
| 500 | внутренняя ошибка                                     |

> **Пока не реализовано:** `gateway` ещё не маппит gRPC-статусы `journal` в HTTP. Сейчас
> любая ошибка `journal` (в том числе доменная: невалидный переход, bullet не найден)
> возвращается как `500` с `"message": "Internal Server Error"`. Если недоступен `identity`,
> middleware сессии отвечает `500` с `"message": "internal server error"` (в нижнем регистре).
> Коды `400`/`404` и сообщения в примерах ниже — целевое поведение. Уже сейчас работают `400`
> на невалидный `id` в пути (`"Bad Request"`) и на тело запроса, которое не разбирается
> `protojson` (`"invalid request body"`).

---

## 3. Эндпоинты

### Health

**GET /health** → `200 { "status": "ok" }` — для Docker/orchestrator healthcheck. Состояние
`identity` и `journal` не проверяется, резолв сессии на этот путь не распространяется.

### Bullets

#### POST /api/bullets

**Юзкейс**: пользователь записывает новую задачу в сегодняшний список
— как строчку в бумажном bullet journal.

- **Актор**: пользователь
- **Предусловие**: нет (bullet создаётся с нуля)
- **Основной поток**: `content` и `kind` не пустые → создаётся bullet
  со `signifier = open`
- **Ошибка**: тело не разбирается как `CreateBulletRequest` (в том числе
  неизвестное значение `kind`, например `"task"`) → `400`
  (`"invalid request body"`); `content` пустой или длиннее 200 символов,
  либо `kind` не указан → `400` (`"bullet content is required"`,
  `"bullet content is too long"`, `"bullet kind must be task, event or note"`)

```json
// request
{
  "content": "Оплатить хостинг",
  "kind": "KIND_TASK"
}
```

```json
// response 201
{
  "data": {
    "id": "uuid",
    "kind": "KIND_TASK",
    "signifier": "SIGNIFIER_OPEN",
    "content": "Оплатить хостинг",
    "userId": "uuid",
    "createdAt": "2026-09-07T08:20:51.871177Z",
    "updatedAt": "2026-09-07T08:20:51.871177Z"
  }
}
```

```json
// response 400
{
  "data": null,
  "error": {
    "code": 400,
    "message": "bullet content is required"
  }
}
```

#### GET /api/bullets

**Юзкейс**: пользователь просматривает свои задачи, сгруппированные
по дням — как разворот страниц в bullet journal.

- **Актор**: пользователь
- **Предусловие**: нет
- **Основной поток**: возвращается плоский список всех bullets
  пользователя
- **Альтернативный поток**: у пользователя ещё нет ни одного bullet →
  `journal` в одной транзакции создаёт 3 demo-bullet'а (с `createdAt`
  на сутки раньше, чтобы их можно было смигрировать) и возвращает их

Группировки по дням в ответе нет: её делает фронтенд
(`web/src/pages/journal/api/bulletsApi.ts`, по локальной дате `createdAt`).

Если список всё же пуст, `protojson` опускает пустое repeated-поле, и
ответ приходит как `{}`, без ключа `data`.

```json
// response 200
{
  "data": [
    /* Bullet[] */
  ]
}
```

#### POST /api/bullets/{id}/complete

**Юзкейс**: пользователь отмечает задачу выполненной, не удаляя её —
как галочка в бумажном bullet journal.

- **Актор**: пользователь
- **Предусловие**: bullet с данным `id` существует, имеет `kind = task`
  и `signifier = open`
- **Основной поток**: bullet переводится в `signifier = completed`
- **Ошибка**: bullet с данным `id` не найден → `404`; bullet уже
  `completed` → `400` (`"bullet already completed"`); bullet в другом
  закрытом состоянии → `400` (`"bullet is not open"`); bullet не задача →
  `400` (`"bullet is not a task"`)

Тело запроса отсутствует. Запись не удаляется, меняется только
`signifier`. Операция не идемпотентна: повторный вызов возвращает `400`.

```json
// response 200
{
  "data": {
    "id": "uuid",
    "kind": "KIND_TASK",
    "signifier": "SIGNIFIER_COMPLETED",
    "content": "Оплатить хостинг",
    "userId": "uuid",
    "createdAt": "2026-09-07T08:20:51.871177Z",
    "updatedAt": "2026-09-07T08:23:20.751189209Z"
  }
}
```

```json
// response 404
{
  "data": null,
  "error": {
    "code": 404,
    "message": "bullet not found"
  }
}
```

#### POST /api/bullets/{id}/migrate

**Юзкейс**: пользователь переносит незавершённую задачу на текущий
день. Задача целиком переезжает в текущий день, а в исходном дне от
неё ничего не остаётся.

- **Актор**: пользователь
- **Предусловие**: bullet существует, имеет `kind = task`,
  `signifier = open` и создан не сегодня
- **Основной поток**: тот же bullet (тот же `id`, `content`,
  `signifier = open`) получает `createdAt` и `updatedAt`, равные
  моменту переноса, и попадает в текущий день
- **Ошибка**: bullet не найден → `404`; bullet не является открытой
  задачей (`kind != task` или `signifier != open`) → `400`
  (`"bullet is not an open task"`); bullet создан сегодня → `400`
  (`"bullet was created today"`)

Тело запроса отсутствует. Перенос всегда выполняется на текущий день.
Новая запись не создаётся, поэтому исходная дата задачи теряется. В
ответе возвращается перенесённый bullet.

```json
// response 200
{
  "data": {
    "id": "uuid",
    "kind": "KIND_TASK",
    "signifier": "SIGNIFIER_OPEN",
    "content": "Оплатить хостинг",
    "userId": "uuid",
    "createdAt": "2026-09-08T09:00:00.000000Z",
    "updatedAt": "2026-09-08T09:00:00.000000Z"
  }
}
```

```json
// response 400 — bullet не открытая задача
{
  "data": null,
  "error": {
    "code": 400,
    "message": "bullet is not an open task"
  }
}
```

```json
// response 400 — bullet создан сегодня
{
  "data": null,
  "error": {
    "code": 400,
    "message": "bullet was created today"
  }
}
```

#### POST /api/bullets/{id}/cancel

**Юзкейс**: пользователь отменяет запись, которая больше не актуальна,
не удаляя её — как зачёркнутая строка в бумажном bullet journal.

- **Актор**: пользователь
- **Предусловие**: bullet с данным `id` существует, имеет
  `signifier = open` (любой `kind`, не только `task`)
- **Основной поток**: bullet переводится в `signifier = cancelled`
- **Ошибка**: bullet с данным `id` не найден → `404`; bullet уже
  `cancelled` → `400` (`"bullet already cancelled"`); bullet в другом
  закрытом состоянии → `400` (`"bullet is not open"`)

Тело запроса отсутствует. Операция не идемпотентна: повторный вызов
возвращает `400`.

```json
// response 200
{
  "data": {
    "id": "bac1e168-356d-41f5-ab1b-c5c1046a6158",
    "kind": "KIND_TASK",
    "signifier": "SIGNIFIER_CANCELLED",
    "content": "Оплатить хостинг",
    "userId": "2f168d69-2bd3-4691-ad7f-72addb4270e0",
    "createdAt": "2026-09-07T08:20:51.871177Z",
    "updatedAt": "2026-09-07T08:23:20.751189209Z"
  }
}
```

```json
// response 404
{
  "data": null,
  "error": {
    "code": 404,
    "message": "bullet not found"
  }
}
```

```json
// response 400
{
  "data": null,
  "error": {
    "code": 400,
    "message": "bullet is not open"
  }
}
```

### Collections

> Реализация отложена: сначала MVP на `Bullet`, `Collection`
> подключается позже.

#### POST /api/collections

**Юзкейс**: пользователь заводит новую коллекцию — тематический
список, в который потом будут попадать bullets, как открыть новый
раздел в bullet journal.

- **Актор**: пользователь
- **Предусловие**: нет (коллекция создаётся с нуля)
- **Основной поток**: `topic` не пустой → создаётся коллекция
- **Ошибка**: `topic` пустой/невалидный → `400`

```json
// request
{ "topic": "Идеи для проекта" }
```

```json
// response 201
{
  "data": {
    "id": "uuid",
    "topic": "Идеи для проекта",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

#### GET /api/collections

**Юзкейс**: пользователь просматривает список своих коллекций, чтобы
выбрать, в какую добавлять bullets.

- **Актор**: пользователь
- **Предусловие**: нет (список может быть пустым)
- **Основной поток**: возвращается список коллекций
- **Альтернативный поток**: коллекций нет → `200` с пустым списком
  (`"data": []`), не `null`

```json
// response 200
{
  "data": [
    /* Collection[] */
  ]
}
```

#### GET /api/collections/{id}

**Юзкейс**: пользователь открывает конкретную коллекцию.

- **Актор**: пользователь
- **Предусловие**: коллекция с данным `id` существует
- **Основной поток**: возвращается коллекция
- **Ошибка**: коллекция с данным `id` не найдена → `404`

```json
// response 200
{
  "data": {
    "id": "uuid",
    "topic": "Идеи для проекта",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

```json
// response 404
{
  "data": null,
  "error": {
    "code": 404,
    "message": "collection not found"
  }
}
```

#### PUT /api/collections/{id}

**Юзкейс**: пользователь переименовывает коллекцию.

- **Актор**: пользователь
- **Предусловие**: коллекция с данным `id` существует
- **Основной поток**: `topic` не пустой → коллекция обновляется
- **Ошибка**: коллекция не найдена → `404`; `topic` пустой/невалидный
  → `400`

```json
// request — полная замена
{ "topic": "Идеи для проекта" }
```

```json
// response 200
{
  "data": {
    "id": "uuid",
    "topic": "Идеи для проекта",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

```json
// response 404
{
  "data": null,
  "error": {
    "code": 404,
    "message": "collection not found"
  }
}
```

```json
// response 400
{
  "data": null,
  "error": {
    "code": 400,
    "message": "topic is required"
  }
}
```

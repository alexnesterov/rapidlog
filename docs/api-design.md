# API — проектирование

Слоистая архитектура: `handler → service → repository (interface)`.
PostgreSQL (`pgx`) — единственное место записи (service layer). RabbitMQ —
side-effects (уведомления, аудит), не основной путь записи.

---

## 1. Сущности

### Bullet

| Поле | Тип | Описание |
| --- | --- | --- |
| id | UUID | |
| collection_id | UUID | коллекция-владелец (FK → Collection) |
| title | string | обязательное, ≤200 символов |
| status | string | `OPEN` \| `DONE` |
| created_at | timestamp | |
| updated_at | timestamp | |

**Операции:** Create, List, Get, Update, Delete, Done (пометить
выполненной; запись не удаляется — удаление только через отдельный
Delete по решению пользователя).

### Collection

> Реализация отложена: сначала MVP на `Bullet`, `Collection`
> подключается позже.

По аналогии с Bullet Journal: каждый bullet живёт в какой-то
коллекции — произвольном именованном списке.

| Поле | Тип | Описание |
| --- | --- | --- |
| id | UUID | |
| topic | string | обязательное, ≤100 символов |
| created_at | timestamp | |
| updated_at | timestamp | |

**Операции:** Create, List, Get, Update, Delete.

---

## 2. Общий формат ответов

**Успех — одиночный объект:**

```json
{ "data": { /* сущность */ } }
```

**Успех — список:**

```json
{ "data": [ /* Bullet[] */ ] }
```

**Ошибка:**

```json
{
  "data": null,
  "error": {
    "code": 400,
    "message": "title is required"
  }
}
```

`error.code` — числовой HTTP-статус (дублирует статус-код ответа),
отдельного строкового кода ошибки нет.

| Код | Когда |
| --- | --- |
| 200 | успешный GET/PUT/POST без создания ресурса |
| 201 | ресурс создан |
| 204 | успешное удаление, тело пустое |
| 400 | некорректный запрос (невалидный JSON) |
| 404 | ресурс не найден |
| 409 | конфликт (удаление непустой коллекции) |
| 422 | семантическая ошибка валидации |
| 500 | внутренняя ошибка |

---

## 3. Эндпоинты

### Health

**GET /health** → `200 { "status": "ok" }` — для Docker/orchestrator
healthcheck.

### Bullets

#### POST /api/bullets

```json
// request
{
  "collection_id": "uuid", "title": "Оплатить хостинг"
}
```

```json
// response 201
{
  "data": {
    "id": "uuid",
    "collection_id": "uuid",
    "title": "Оплатить хостинг",
    "status": "OPEN",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

Побочный эффект: публикуется событие `bullet.created` (см. раздел 4).

#### GET /api/bullets?collection_id=

- `collection_id` — фильтр по коллекции (опционально)

Ответ группирует bullets по календарному дню (`created_at`, дата без
времени): группы отсортированы от нового дня к старому, а bullets
внутри группы — от старого к новому (новые снизу).

```json
// response 200
{
  "data": [
    {
      "day": "2026-07-28",
      "bullets": [ /* Bullet[] */ ]
    }
  ]
}
```

#### GET /api/bullets/{id}

Ответ 200 — Bullet, либо `404`.

#### PUT /api/bullets/{id}

```json
// request — полная замена
{
  "collection_id": "uuid", "title": "...",
  "status": "OPEN"
}
```

Ответ 200 — обновлённый Bullet. `PUT` — это полная замена состояния,
поэтому `status` обязателен в теле: смена `collection_id` переносит
bullet в другую коллекцию (миграция на завтра), а смена `status` с
`DONE` на `OPEN` — это и есть reopen, отдельного эндпоинта для него
нет. В отличие от `POST /done`, `PUT` не публикует `bullet.completed`
— это чисто CRUD-замена, без бизнес-события.

#### DELETE /api/bullets/{id}

Ответ `204`, либо `404` — если bullet с таким `id` не существует.
Публикует `bullet.deleted` только при реальном удалении.

#### POST /api/bullets/{id}/done

Логика: `status` bullet переводится в `DONE`, запись не удаляется —
удалить её пользователь может отдельно через `DELETE`. Идемпотентно:
если bullet уже `DONE`, повторный вызов — `200` без изменений и без
повторной публикации `bullet.completed`.

```json
// response 200
{
  "data": {
    "id": "uuid",
    "title": "Оплатить хостинг",
    "status": "DONE",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

Публикует `bullet.completed`.

### Collections

> Реализация отложена: сначала MVP на `Bullet`, `Collection`
> подключается позже.

#### POST /api/collections

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

```json
// response 200
{ "data": [ /* Collection[] */ ] }
```

#### GET /api/collections/{id}

Ответ 200 — Collection, либо `404`.

#### PUT /api/collections/{id}

```json
// request — полная замена
{ "topic": "Идеи для проекта" }
```

Ответ 200 — обновлённая Collection.

#### DELETE /api/collections/{id}

Ответ `204`, либо `404` — если коллекция не существует. Если внутри
есть bullets — `409`, сначала перенеси или удали их. Источник истины —
`FK bullets.collection_id REFERENCES collections(id) ON DELETE
RESTRICT` в схеме Postgres, а не
предварительная проверка в service-слое: под конкурентной нагрузкой
check-then-delete даёт гонку (bullet может быть создан между проверкой
и удалением). `409` в этом эндпоинте — это маппинг ошибки нарушения
FK-констрейнта (`23503`) на HTTP-код, а не отдельная бизнес-проверка.

---

## 4. Асинхронные события (RabbitMQ)

Exchange: `bullet_events` (topic).

- **`bullet.created`** — после успешного `POST /api/bullets`.
  Payload: `{bullet_id, collection_id, title}`.
- **`bullet.completed`** — после `POST /api/bullets/{id}/done`.
  Payload: `{bullet_id, collection_id, title, completed_at}`.
- **`bullet.deleted`** — после `DELETE /api/bullets/{id}`.
  Payload: `{bullet_id, collection_id, deleted_at}`.

Consumer (планируется отдельным воркером): отправка уведомлений, запись в
audit-log. Service слой публикует событие уже после успешной записи в
Postgres — паттерн "transactional write + async notification" (без
outbox на этом этапе, риск dual-write принят как допустимый для
pet-проекта).

---

## 5. Поток создания bullet (кратко)

```text
Client → HTTP handler → service.CreateBullet()
  1. repository.Insert() → Postgres (синхронно, ответ клиенту)
  2. publisher.Publish("bullet.created") → RabbitMQ (best-effort)
```

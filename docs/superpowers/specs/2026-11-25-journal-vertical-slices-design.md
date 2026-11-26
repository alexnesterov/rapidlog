# Journal → вертикальные слайсы (дизайн)

## Context

Цель учебная: пощупать архитектуру вертикальных слайсов на живом сервисе и сравнить со
слоистой схемой. Сейчас `journal` — `controller/grpcapi → usecase ← adapter/postgres` + `entity`;
после недавней серии PR (#127–#131) у каждого юзкейса уже своя структура `Execute(ctx, XxxInput)`
и узкий интерфейс в `grpcapi`, но фича размазана по трём папкам и общему `BulletRepository`.
Внешнее поведение (gRPC-контракт, `gateway`, БД-схема) не меняется. Код пишу я, пользователь
читает и сравнивает.

## Approaches (выбран B)

- **A — файлы по фичам, общий репозиторий.** Механический переезд, но по сути те же слои. Слабый урок.
- **B — слайс владеет своим SQL (выбран).** Общие только доменные правила, маппер и тип
  gRPC-сервера; каждый слайс — хендлер + логика + свой узкий `Store` с собственными запросами.
  Видна суть подхода, домен не дублируется.
- **C — слайс = хендлер, сам ходит в pgx.** Меньше всего кода, но юнит-тесты на моках
  уходят в интеграционные — против принятой в проекте стратегии тестирования.

## Целевая структура

```
journal/internal/
  entity/                 ← без изменений (Bullet, переходы, ValidationError, ErrNotFound)
  protomap/protomap.go    ← бывший grpcapi/mapper.go (ToProtoBullet)
  server/server.go        ← JournalService: Unimplemented… + по полю-интерфейсу на RPC,
                            методы — однострочное делегирование; Register(srv, Handlers, …)
  postgres/               ← только Connect + Migrate (бывший adapter/postgres без репозитория)
  features/
    listbullets/          handler.go, listbullets.go, store.go, *_test.go, mocks/
    createbullet/
    completebullet/
    cancelbullet/
    migratebullet/
```

Внутри слайса (на примере `completebullet`):
- `handler.go` — `Handler{uc Usecase, logger}` с `Handle(ctx, *journalv1.CompleteBulletRequest)`:
  парсинг uuid, вызов `Usecase.Execute`, маппинг ошибок (`*entity.ValidationError` →
  `InvalidArgument`, прочее → `Internal`), `protomap.ToProtoBullet`. Логика = текущий
  `grpcapi/complete_bullet.go`.
- `completebullet.go` — `Input`, интерфейсы `Usecase` и `Store`, `CompleteBullet` с `Execute`
  (логика = текущий `usecase/complete_bullet.go`).
- `store.go` — `PostgresStore{pool}`, только нужные слайсу запросы, через
  `transactor.GetQueryer` (`internal/pkg/transactor`).
- Тесты: `handler_test.go` на `mocks.MockUsecase`, `completebullet_test.go` на `mocks.MockStore`
  (переносятся из текущих тестов). Store без юнит-тестов, как прежние адаптеры.

## SQL по слайсам (вместо общего upsert `Save`)

| Слайс | Store |
|---|---|
| list | `List(userID)` — текущий SELECT |
| create | `Insert(b)` — INSERT без ON CONFLICT |
| complete, cancel | `Get(id, userID)` + `UpdateSignifier(b)`: `SET signifier, updated_at WHERE id AND user_id` |
| migrate | `Get(id, userID)` + `UpdateDates(b)`: `SET created_at, updated_at WHERE id AND user_id` |

0 затронутых строк в UPDATE → `entity.ErrNotFound`; `ErrOwnerMismatch` удаляется.
`SELECT` и Scan сознательно дублируются по слайсам (мои решения: дублирование вместо общего пакета;
`entity` не переименовывать в `domain`, чтобы не раздувать дифф).

## Порядок (одна ветка `refactor/journal-vertical-slices` от актуального `main`, один PR;
## шаги идут по очереди, после каждого сервис собирается и тесты зелёные)

0. `protomap` (вынос маппера), запись `internal/service/journal/internal/features` в
   `.mockery.yml` (`all: true`, `recursive: true`).
1. `listbullets` → `grpcapi.JournalService.ListBullets` делегирует в слайс.
2. `createbullet`, 3. `completebullet`, 4. `cancelbullet`, 5. `migratebullet` — тем же образом.
6. Уборка: `controller/grpcapi` → `server` (только делегирование, без тестов), удалить пакет
   `usecase`, `adapter/postgres.BulletRepository` (→ `postgres` с Connect/Migrate),
   `ErrOwnerMismatch`, лишние записи в `.mockery.yml`, старые `mocks/`; сборка в `app.go`.
   Обновить CLAUDE.md (раздел journal, тестирование, mockery) и README (структура).

Всё в одной ветке. Коммитить — только по просьбе пользователя (`/git-branch-commit`);
предлагаемые сабджекты по шагам: `refactor(journal): move <x> to vertical slice`.

## Verification

- После каждого шага: `task mock`, `gofmt -l`, `go vet ./internal/service/journal/...`,
  `go test ./internal/service/journal/...` (запускаю сам — статические проверки разрешены).
- Живую проверку делает пользователь (не запускаю сервисы): `docker compose up -d --build`,
  затем Bruno `journal/*.bru` или `bullets/*.bru` — create → list → complete/cancel/migrate,
  плюс чужой `user_id` на complete → ошибка (раньше OwnerMismatch/NotFound, теперь NotFound).

## После одобрения

Сохранить этот дизайн в `docs/superpowers/specs/2026-11-25-journal-vertical-slices-design.md`
(коммитом), затем `superpowers:writing-plans` для пошагового плана.

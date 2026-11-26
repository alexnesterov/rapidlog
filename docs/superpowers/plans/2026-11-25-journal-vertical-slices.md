# Journal Vertical Slices Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Переложить `journal` со слоёв (`controller/grpcapi → usecase ← adapter/postgres`) на вертикальные слайсы, где каждая RPC-фича владеет хендлером, логикой и своим SQL.

**Architecture:** Пять пакетов `internal/features/<rpc>` (handler + usecase + postgres store), общие только `entity`, `protomap` и `server` (gRPC-тип с однострочным делегированием). Переезд по одной RPC; после каждой задачи сервис собирается и тесты зелёные.

**Tech Stack:** Go, gRPC (`journalv1`), pgx v5 (`pgxpool`), `internal/pkg/transactor.GetQueryer`, testify (suite / табличные тесты), mockery v3 (`task mock`).

**Spec:** `docs/superpowers/specs/2026-11-25-journal-vertical-slices-design.md`

## Global Constraints

- Всё в ветке `refactor/journal-vertical-slices`; **не коммитить**, пока пользователь не попросит (тогда `/git-branch-commit`, сабджект `refactor(journal): move <x> to vertical slice`).
- gRPC-контракт, `gateway`, схема БД и миграции не меняются.
- `entity` не переименовывать; логика юзкейсов и маппинг ошибок хендлеров переносятся как есть (`*entity.ValidationError` → `codes.InvalidArgument` с текстом ошибки, прочее → `codes.Internal, "internal error"`).
- Postgres-store без юнит-тестов (как прежние адаптеры); сервисы и docker не запускать — живую проверку делает пользователь.
- После `task mock` вручную удалять устаревшие `mock_*.go` (mockery их не чистит).
- Проверка каждой задачи: `task mock && gofmt -l internal/service/journal && go vet ./internal/service/journal/... && go test ./internal/service/journal/...` → пустой `gofmt`, все `ok`.

## Слайс: имена (одинаковы во всех пяти пакетах)

В пакете `features/<name>` (`listbullets`, `createbullet`, `completebullet`, `cancelbullet`, `migratebullet`):

| Файл | Содержимое |
|---|---|
| `<name>.go` | `type Input struct{…}`; `type Store interface{…}`; `type Usecase struct{store Store}`; `func New(store Store) *Usecase`; `func (uc *Usecase) Execute(ctx context.Context, in Input) (<result>, error)` |
| `handler.go` | `type Executor interface{ Execute(ctx context.Context, in Input) (<result>, error) }`; `type Handler struct{uc Executor; logger *slog.Logger}`; `func NewHandler(uc Executor, logger *slog.Logger) *Handler`; `func (h *Handler) Handle(ctx context.Context, req *journalv1.<Rpc>Request) (*journalv1.<Rpc>Response, error)` |
| `store.go` | `type PostgresStore struct{pool *pgxpool.Pool}`; `func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore`; методы `Store`; запросы через `transactor.GetQueryer(ctx, s.pool)`; `pgx.ErrNoRows` → `entity.ErrNotFound` |
| `<name>_test.go` | тест `Usecase` на `mocks.MockStore` (пакет `<name>_test`) |
| `handler_test.go` | тест `Handler` на `mocks.MockExecutor` (пакет `<name>_test`) |

`<result>` = `*entity.Bullet`, у `listbullets` — `[]*entity.Bullet`. `Input`: `listbullets{UserID}`, `createbullet{UserID, Kind, Content}`, остальные `{ID, UserID}`.

## Review Focus

1. Чужой `user_id` на complete/cancel/migrate — `UPDATE … WHERE id AND user_id` затрагивает 0 строк → `entity.ErrNotFound` (не тихий успех). Тест: кейс `update not found` в usecase-тестах complete/cancel/migrate (store возвращает `ErrNotFound` → `Execute` возвращает его, bullet `nil`).
2. Пустой журнал нового пользователя — list отдаёт пустой список без ошибки (кейс `Empty` переносится).
3. Повторный complete/cancel идемпотентен и **не** пишет в БД (кейсы `AlreadyCompleted`/`AlreadyCanceled` без ожидания `UpdateSignifier`; mockery-мок упадёт на неожиданном вызове).
4. Migrate не-task / не-open / созданного сегодня — доменная ошибка, `UpdateDates` не вызывается (кейс `NotOpenTask` переносится без ожидания записи).
5. Create не перезаписывает существующую запись: `INSERT` без `ON CONFLICT`, конфликт id → ошибка pgx → `Internal` (store без тестов — проверяется ревью SQL и вручную).

---

### Task 0: protomap и mockery для слайсов

**Files:**
- Create: `internal/service/journal/internal/protomap/protomap.go` (переезд `controller/grpcapi/mapper.go`, `package protomap`)
- Delete: `internal/service/journal/internal/controller/grpcapi/mapper.go`
- Modify: все `controller/grpcapi/*.go` и `*_test.go`, где `ToProtoBullet`/`ToProtoBullets` → `protomap.…`
- Modify: `.mockery.yml` — добавить пакет `github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features` с `config: {all: true, recursive: true}`

**Interfaces:**
- Produces: `protomap.ToProtoKind(entity.Kind) (journalv1.Kind, error)`, `protomap.ToProtoSignifier(entity.Signifier) (journalv1.Signifier, error)`, `protomap.ToProtoBullet(*entity.Bullet) (*journalv1.Bullet, error)`, `protomap.ToProtoBullets([]*entity.Bullet) ([]*journalv1.Bullet, error)`

- [ ] **Step 1:** Перенести файл, сменить пакет, обновить вызовы в `grpcapi` и его тестах.
- [ ] **Step 2:** Добавить запись в `.mockery.yml`. Пакета `features` ещё нет — `task mock` не должен падать (если падает на отсутствующем пакете, добавить запись в Task 1 вместе с первым слайсом).
- [ ] **Step 3:** Проверка (см. Global Constraints) → все `ok`.

### Task 1: listbullets

**Files:**
- Create: `internal/service/journal/internal/features/listbullets/{listbullets.go,handler.go,store.go,listbullets_test.go,handler_test.go}`
- Delete: `controller/grpcapi/list_bullets.go`, `controller/grpcapi/list_bullets_test.go`, `usecase/list_bullets.go`, `usecase/list_bullets_test.go`
- Modify: `controller/grpcapi/service.go`, `internal/service/journal/app.go`

**Interfaces:**
- Consumes: `protomap.ToProtoBullets` (Task 0)
- Produces: `listbullets.Store{ List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) }`; `listbullets.New`, `listbullets.NewHandler`, `listbullets.NewPostgresStore` (таблица имён выше).
- Produces (переходный шов в `grpcapi/service.go`, расширяется в Tasks 2–5): `type ListBulletsHandler interface{ Handle(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error) }`; `type Handlers struct{ ListBullets ListBulletsHandler }`; `NewJournalService(uc Usecases, h Handlers, logger *slog.Logger)`, `RegisterJournalService(srv *grpc.Server, uc Usecases, h Handlers, logger *slog.Logger)`; метод `JournalService.ListBullets` — `return s.h.ListBullets.Handle(ctx, req)`. Поле `Usecases.ListBullets` и интерфейс `ListBulletsUsecase` удаляются.

- [ ] **Step 1: Перенести тесты.** `usecase/list_bullets_test.go` → `listbullets_test.go` (suite на `mocks.NewMockStore`, `listbullets.New(store)`, `Execute(ctx, listbullets.Input{UserID: …})`; кейсы `ListError`, `Empty`, `Success` с точным `userID`). `grpcapi/list_bullets_test.go` → `handler_test.go` (`mocks.NewMockExecutor`, `listbullets.NewHandler(m, slog.Default()).Handle(...)`, ожидание `Execute(mock.Anything, listbullets.Input{UserID: userID})`; кейсы как сейчас).
- [ ] **Step 2:** `go test ./internal/service/journal/internal/features/...` → FAIL компиляции (`undefined: listbullets.New` и т.п.).
- [ ] **Step 3:** Реализовать `listbullets.go` (логика = `usecase/list_bullets.go`), `handler.go` (= `grpcapi/list_bullets.go`), `store.go` (SELECT + Scan = `adapter/postgres.BulletRepository.List`). Удалить старые файлы, перевести `grpcapi` на `Handlers`, в `app.go`: `listbullets.NewHandler(listbullets.New(listbullets.NewPostgresStore(pool)), logger)`. Остальные тесты `grpcapi` передают `grpcapi.Handlers{}`.
- [ ] **Step 4:** Проверка → все `ok`; в `features/listbullets/mocks/` есть `MockStore` и `MockExecutor`; `mocks.MockListBulletsUsecase` исчез из `grpcapi/mocks`.

### Task 2: createbullet

**Files:** как в Task 1 для `create_bullet` → `features/createbullet/`.

**Interfaces:**
- Produces: `createbullet.Store{ Insert(ctx context.Context, b *entity.Bullet) error }`; `Input{UserID uuid.UUID; Kind entity.Kind; Content string}`; `Handlers.CreateBullet CreateBulletHandler` (Handle с `*journalv1.CreateBulletRequest`/`Response`).

- [ ] **Step 1: Перенести тесты.** `usecase/create_bullet_test.go` → `createbullet_test.go` (ожидание `Insert` вместо `Save`); `grpcapi/create_bullet_test.go` → `handler_test.go` (включая кейс `parse kind error` → `InvalidArgument "invalid kind"`, мок пустой).
- [ ] **Step 2:** Тесты слайса → FAIL компиляции.
- [ ] **Step 3:** Реализовать. Хендлер сохраняет `entity.ParseKind(protoenum.ToCleanString(req.GetKind()))`. `Insert`: `INSERT INTO bullets (id, kind, signifier, content, user_id, created_at, updated_at) VALUES ($1..$7)` — без `ON CONFLICT`. Удалить старые файлы, `Usecases.CreateBullet`, обновить `app.go`.
- [ ] **Step 4:** Проверка → все `ok`.

### Task 3: completebullet

**Interfaces:**
- Produces: `completebullet.Store{ Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error); UpdateSignifier(ctx context.Context, b *entity.Bullet) error }`; `Handlers.CompleteBullet`.

- [ ] **Step 1: Перенести тесты** (`Save` → `UpdateSignifier`, имена кейсов `TestCompleteBullet_*`) и добавить кейс:

```go
func (s *CompleteBulletSuite) TestCompleteBullet_UpdateNotFound() {
	in := completebullet.Input{ID: uuid.New(), UserID: uuid.New()}
	s.store.EXPECT().Get(mock.Anything, in.ID, in.UserID).
		Return(&entity.Bullet{Kind: entity.KindTask, Signifier: entity.SignifierOpen}, nil).Once()
	s.store.EXPECT().UpdateSignifier(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(entity.ErrNotFound).Once()

	got, err := s.uc.Execute(s.T().Context(), in)
	s.Nil(got)
	s.ErrorIs(err, entity.ErrNotFound)
}
```

- [ ] **Step 2:** → FAIL компиляции.
- [ ] **Step 3:** Реализовать. `Get` = прежний SELECT по `id AND user_id`. `UpdateSignifier`: `UPDATE bullets SET signifier = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`; `RowsAffected() == 0` → `entity.ErrNotFound`. Удалить старые файлы, `Usecases.CompleteBullet`, обновить `app.go`.
- [ ] **Step 4:** Проверка → все `ok`.

### Task 4: cancelbullet

Как Task 3 (своя копия `Get`/`UpdateSignifier` в `features/cancelbullet/store.go` — дублирование сознательное). Тесты из `usecase/cancel_bullet_test.go` и `grpcapi/cancel_bullet_test.go` + кейс `TestCancelBullet_UpdateNotFound` (Get возвращает open-bullet любого kind).

- [ ] **Step 1:** Перенести тесты + новый кейс. **Step 2:** FAIL. **Step 3:** Реализовать, удалить старое, `app.go`. **Step 4:** Проверка → `ok`.

### Task 5: migratebullet

**Interfaces:**
- Produces: `migratebullet.Store{ Get(...) (*entity.Bullet, error); UpdateDates(ctx context.Context, b *entity.Bullet) error }`; `Handlers.MigrateBullet`.

- [ ] **Step 1:** Перенести тесты (`Save` → `UpdateDates`) + кейс `TestMigrateBullet_UpdateNotFound` (bullet: open task, `CreatedAt` сутки назад).
- [ ] **Step 2:** → FAIL.
- [ ] **Step 3:** Реализовать. `UpdateDates`: `UPDATE bullets SET created_at = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`; 0 строк → `ErrNotFound`. Удалить старое; после этого `Usecases` пуст — убрать его и параметр `uc` из `NewJournalService`/`RegisterJournalService`.
- [ ] **Step 4:** Проверка → `ok`; пакет `usecase` содержит только `ports.go` и `mocks/`.

### Task 6: уборка и документация

**Files:**
- Move: `controller/grpcapi/service.go` → `internal/server/server.go` (`package server`; `type JournalService`, `type Handlers`, интерфейсы `XxxHandler`; `func Register(srv *grpc.Server, h Handlers)` — logger больше не нужен, т.к. логируют хендлеры слайсов). Удалить `controller/` целиком (включая `mocks/`).
- Delete: пакет `usecase` (`ports.go`, `mocks/`); `adapter/postgres/bullet_repository.go`; `entity.ErrOwnerMismatch`.
- Move: `adapter/postgres/{postgres.go,migrate.go}` → `internal/postgres/` (`package postgres`, `Connect`, `Migrate`).
- Modify: `app.go`, `.mockery.yml` (убрать записи `journal/internal/usecase` и `journal/internal/controller/grpcapi`).
- Modify: `CLAUDE.md` (раздел journal: структура, слайсы, тестирование, список пакетов mockery; убрать `ErrOwnerMismatch`, upsert `Save`) — вне git; `README.md` (блок «Структура», абзац про mockery).

- [ ] **Step 1:** Перемещения/удаления, `app.go` собирает `server.Register(grpcServer, server.Handlers{…})`.
- [ ] **Step 2:** `grep -rn "grpcapi\|usecase\.\|adapter/postgres\|ErrOwnerMismatch" internal/service/journal .mockery.yml` → пусто.
- [ ] **Step 3:** Проверка → все `ok`; `go build ./...` (нужен `web/dist`) → ok.
- [ ] **Step 4:** Обновить CLAUDE.md и README.
- [ ] **Step 5:** Сообщить пользователю шаги живой проверки: `docker compose up -d --build`, Bruno `bullets/*.bru`: create → list → complete → повторный complete → cancel другого → migrate вчерашнего; complete с чужим токеном/`user_id` → 500 (NotFound → Internal).

package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKeyTx struct{}

type queryer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func getQueryer(ctx context.Context, pool *pgxpool.Pool) queryer {
	if tx, ok := ctx.Value(ctxKeyTx{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}

type transactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *transactionManager {
	return &transactionManager{pool: pool}
}

func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in transaction", "recovered", r, "stack", string(debug.Stack()))
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("panic in transaction: %v", r)
		} else {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				slog.Error("rollback transaction", "error", rollbackErr)
			}
		}
	}()

	ctx = context.WithValue(ctx, ctxKeyTx{}, tx)

	execErr := fn(ctx)
	if execErr != nil {
		return fmt.Errorf("execute transaction: %w", execErr)
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		return fmt.Errorf("commit transaction: %w", commitErr)
	}

	return nil
}

var _ port.TransactionManager = &transactionManager{}

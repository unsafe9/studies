package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unsafe9/studies/go-sqlc/db"
)

type txState struct {
	onCommitSuccessCallbacks []func(context.Context) error
}

type dbCtxKeyType int

var (
	dbCtxKey      = dbCtxKeyType(0)
	txStateCtxKey = dbCtxKeyType(1)
)

func WithDB(ctx context.Context, pool *pgxpool.Pool) context.Context {
	return context.WithValue(ctx, dbCtxKey, pool)
}

func DB(ctx context.Context) db.DBTX {
	return ctx.Value(dbCtxKey).(db.DBTX)
}

func OnCommitSuccess(ctx context.Context, callbacks ...func(context.Context) error) error {
	tx := ctx.Value(txStateCtxKey)
	if tx == nil {
		// execute directly outside of transaction
		for _, callback := range callbacks {
			if err := callback(ctx); err != nil {
				return err
			}
		}
	}

	txs := tx.(*txState)
	txs.onCommitSuccessCallbacks = append(txs.onCommitSuccessCallbacks, callbacks...)
	return nil
}

func Transaction(ctx context.Context, callback func(context.Context) error) error {
	if tx := ctx.Value(txStateCtxKey); tx != nil {
		// Just call if the transaction is already in progress.
		//  pgx supports nested tx based on db savepoint, but we don't use it for now.
		return callback(ctx)
	}

	pool, ok := DB(ctx).(*pgxpool.Pool)
	if !ok {
		return fmt.Errorf("invalid db type: %T", DB(ctx))
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	txs := &txState{}
	txCtx := context.WithValue(ctx, txStateCtxKey, txs)
	txCtx = context.WithValue(txCtx, dbCtxKey, tx)

	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback(txCtx)
			panic(v)
		}
	}()

	if err := callback(txCtx); err != nil {
		if rollbackErr := tx.Rollback(txCtx); rollbackErr != nil {
			err = fmt.Errorf("%w, failed to rollback transaction: %w", err, rollbackErr)
		}
		return err
	}

	if err := tx.Commit(txCtx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	for _, onSuccess := range txs.onCommitSuccessCallbacks {
		// use original context for onSuccess callbacks
		if err := onSuccess(ctx); err != nil {
			return err
		}
	}
	return nil
}

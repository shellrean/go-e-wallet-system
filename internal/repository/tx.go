package repository

import (
	"context"
	"github.com/doug-martin/goqu/v9"
)

type txContextKey struct{}

type dbExecutor interface {
	From(from ...interface{}) *goqu.SelectDataset
	Insert(table interface{}) *goqu.InsertDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
}

func getExecutor(ctx context.Context, db *goqu.Database) dbExecutor {
	if tx, ok := ctx.Value(txContextKey{}).(*goqu.TxDatabase); ok && tx != nil {
		return tx
	}
	return db
}

func isTx(ctx context.Context) bool {
	tx, ok := ctx.Value(txContextKey{}).(*goqu.TxDatabase)
	return ok && tx != nil
}

func withTransaction(ctx context.Context, db *goqu.Database, fn func(ctx context.Context) error) error {
	if isTx(ctx) {
		return fn(ctx)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txContextKey{}, tx)
	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

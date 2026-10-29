package repository

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"shellrean.id/belajar-auth/domain"
)

type transactionRepository struct {
	db *goqu.Database
}

func NewTransaction(con *sql.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: goqu.New("default", con),
	}
}

func (t *transactionRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return withTransaction(ctx, t.db, fn)
}

func (t transactionRepository) Insert(ctx context.Context, transaction *domain.Transaction) error {
	executor := getExecutor(ctx, t.db).Insert("transactions").Rows(goqu.Record{
		"account_id":           transaction.AccountId,
		"sof_number":           transaction.SofNumber,
		"dof_number":           transaction.DofNumber,
		"amount":               transaction.Amount,
		"transaction_type":     transaction.TransactionType,
		"transaction_datetime": transaction.TransactionDatetime,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, transaction)
	return err
}

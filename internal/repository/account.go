package repository

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"shellrean.id/belajar-auth/domain"
)

type accountRepository struct {
	db *goqu.Database
}

func NewAccount(con *sql.DB) domain.AccountRepository {
	return &accountRepository{
		db: goqu.New("default", con),
	}
}

func (a *accountRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return withTransaction(ctx, a.db, fn)
}

func (a accountRepository) FindByUserID(ctx context.Context, id int64) (account domain.Account, err error) {
	dataset := getExecutor(ctx, a.db).From("accounts").Where(goqu.Ex{
		"user_id": id,
	})
	if isTx(ctx) {
		dataset = dataset.ForUpdate(exp.Wait)
	}
	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

func (a accountRepository) FindByAccountNumber(ctx context.Context, accNumber string) (account domain.Account, err error) {
	dataset := getExecutor(ctx, a.db).From("accounts").Where(goqu.Ex{
		"account_number": accNumber,
	})
	if isTx(ctx) {
		dataset = dataset.ForUpdate(exp.Wait)
	}
	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

func (a accountRepository) Insert(ctx context.Context, account *domain.Account) error {
	executor := getExecutor(ctx, a.db).Insert("accounts").Rows(goqu.Record{
		"user_id":        account.UserId,
		"account_number": account.AccountNumber,
		"balance":        account.Balance,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, account)
	return err
}

func (a accountRepository) Update(ctx context.Context, account *domain.Account) error {
	executor := getExecutor(ctx, a.db).Update("accounts").Where(goqu.Ex{
		"id": account.ID,
	}).Set(goqu.Record{
		"balance": account.Balance,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

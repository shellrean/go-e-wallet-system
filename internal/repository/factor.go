package repository

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"shellrean.id/belajar-auth/domain"
)

type factorRepository struct {
	db *goqu.Database
}

func NewFactor(con *sql.DB) domain.FactorRepository {
	return &factorRepository{
		db: goqu.New("default", con),
	}
}

func (f *factorRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return withTransaction(ctx, f.db, fn)
}

func (f factorRepository) FindByUser(ctx context.Context, id int64) (factor domain.Factor, err error) {
	dataset := getExecutor(ctx, f.db).From("factors").Where(goqu.Ex{
		"user_id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &factor)
	return
}

func (f factorRepository) Insert(ctx context.Context, factor *domain.Factor) error {
	executor := getExecutor(ctx, f.db).Insert("factors").Rows(goqu.Record{
		"user_id": factor.UserID,
		"pin":     factor.PIN,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, factor)
	return err
}

package repository

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"shellrean.id/belajar-auth/domain"
)

type templateRepository struct {
	db *goqu.Database
}

func NewTemplate(conn *sql.DB) domain.TemplateRepository {
	return &templateRepository{
		db: goqu.New("default", conn),
	}
}

func (t templateRepository) FindByCode(ctx context.Context, code string) (temp domain.Template, err error) {
	dataset := t.db.From("templates").Where(goqu.Ex{
		"code": code,
	})
	_, err = dataset.ScanStructContext(ctx, &temp)
	return
}

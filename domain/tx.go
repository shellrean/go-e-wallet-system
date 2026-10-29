package domain

import "context"

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func RunInTransaction(ctx context.Context, repo interface{}, fn func(ctx context.Context) error) error {
	if txRepo, ok := repo.(Transactor); ok {
		return txRepo.WithTransaction(ctx, fn)
	}
	return fn(ctx)
}

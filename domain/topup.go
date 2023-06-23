package domain

import (
	"context"
	"shellrean.id/belajar-auth/dto"
)

type TopUp struct {
	ID      string  `db:"id"`
	UserID  int64   `db:"user_id"`
	Status  int8    `db:"status"`
	Amount  float64 `db:"amount"`
	SnapURL string  `db:"snap_url"`
}

type TopUpRepository interface {
	FindById(ctx context.Context, id string) (TopUp, error)
	Insert(ctx context.Context, t *TopUp) error
	Update(ctx context.Context, t *TopUp) error
}

type TopUpService interface {
	ConfirmedTopUp(ctx context.Context, id string) error
	InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error)
}

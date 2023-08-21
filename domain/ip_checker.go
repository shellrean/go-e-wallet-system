package domain

import (
	"context"
	"shellrean.id/belajar-auth/dto"
)

type IpCheckerService interface {
	Query(ctx context.Context, ip string) (dto.IpChecker, error)
}

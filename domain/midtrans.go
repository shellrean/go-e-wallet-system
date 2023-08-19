package domain

import "context"

type MidtransService interface {
	GenerateSnapURL(ctx context.Context, t *TopUp) error
	VerifyPayment(ctx context.Context, orderId string) (bool, error)
}

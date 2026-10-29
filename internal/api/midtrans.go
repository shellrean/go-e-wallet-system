package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
)

type midtransApi struct {
	midtransService domain.MidtransService
	topUpService    domain.TopUpService
}

func NewMidtrans(app *fiber.App,
	midtransService domain.MidtransService,
	topUpService domain.TopUpService) {
	m := midtransApi{
		midtransService: midtransService,
		topUpService:    topUpService,
	}

	app.Post("/midtrans/payment-callback", m.paymentHandlerNotification)
}

func (m midtransApi) paymentHandlerNotification(ctx *fiber.Ctx) error {
	var notificationPayload map[string]interface{}
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return ctx.SendStatus(400)
	}

	orderId, exists := notificationPayload["order_id"].(string)
	if !exists || orderId == "" {
		return ctx.SendStatus(400)
	}

	success, err := m.midtransService.VerifyPayment(ctx.Context(), orderId)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}
	if success {
		err = m.topUpService.ConfirmedTopUp(ctx.Context(), orderId)
		if err != nil && !errors.Is(err, domain.ErrTopUpAlreadyConfirmed) {
			return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
				Message: err.Error(),
			})
		}
		return ctx.SendStatus(200)
	}

	return ctx.SendStatus(400)
}

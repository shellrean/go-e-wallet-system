package api

import (
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/domain"
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
	if !exists {
		return ctx.SendStatus(400)
	}

	success, _ := m.midtransService.VerifyPayment(ctx.Context(), notificationPayload)
	if success {
		_ = m.topUpService.ConfirmedTopUp(ctx.Context(), orderId)
		return ctx.SendStatus(200)
	}

	return ctx.SendStatus(400)
}

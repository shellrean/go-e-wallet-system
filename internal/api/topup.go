package api

import (
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
)

type topUpApi struct {
	topUpService domain.TopUpService
}

func NewTopUp(app *fiber.App, authMid fiber.Handler, topService domain.TopUpService) {
	t := topUpApi{
		topUpService: topService,
	}

	app.Post("/topup/initialize", authMid, t.InitializeTopUp)
}

func (t topUpApi) InitializeTopUp(ctx *fiber.Ctx) error {
	var req dto.TopUpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-user").(dto.UserData)
	req.UserID = user.ID

	res, err := t.topUpService.InitializeTopUp(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(res)
}

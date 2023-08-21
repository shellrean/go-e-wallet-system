package api

import (
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
)

type authApi struct {
	userService domain.UserService
	fdsService domain.FdsService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMid fiber.Handler, fdsService domain.FdsService) {
	h := authApi{
		userService: userService,
		fdsService: fdsService,
	}

	app.Post("token/generate", h.GenerateToken)
	app.Get("token/validate", authMid, h.ValidateToken)
	app.Post("user/register", h.RegisterUser)
	app.Post("user/validate-opt", h.ValidateOTP)
}

func (a authApi) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	token, err := a.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}
	if !a.fdsService.IsAuthorized(ctx.Context(), ctx.Get("X-FORWARDED-FOR"), token.UserId) {
		return ctx.SendStatus(401)
	}
	return ctx.Status(200).JSON(token)
}

func (a authApi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")
	return ctx.Status(200).JSON(user)
}

func (a authApi) RegisterUser(ctx *fiber.Ctx) error {
	var req dto.UserRegisterReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	res, err := a.userService.Register(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}
	return ctx.Status(200).JSON(res)
}

func (a authApi) ValidateOTP(ctx *fiber.Ctx) error {
	var req dto.ValidateOtpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	err := a.userService.ValidateOTP(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}
	return ctx.SendStatus(200)
}

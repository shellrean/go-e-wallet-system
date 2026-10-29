package middleware

import (
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/internal/util"
	"strings"
)

func Authenticate(userService domain.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := strings.TrimSpace(ctx.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			authHeader = strings.TrimSpace(authHeader[7:])
		}
		token := authHeader
		if token == "" {
			return ctx.SendStatus(401)
		}
		user, err := userService.ValidateToken(ctx.Context(), token)
		if err != nil {
			return ctx.SendStatus(util.GetHttpStatus(err))
		}
		ctx.Locals("x-user", user)
		return ctx.Next()
	}
}

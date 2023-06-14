package main

import (
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/internal/api"
	"shellrean.id/belajar-auth/internal/component"
	"shellrean.id/belajar-auth/internal/config"
	"shellrean.id/belajar-auth/internal/middleware"
	"shellrean.id/belajar-auth/internal/repository"
	"shellrean.id/belajar-auth/internal/service"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	cacheConnection := component.GetCacheConnection()

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, authMid, transactionService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}

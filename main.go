package main

import (
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/api"
	"shellrean.id/belajar-auth/internal/component"
	"shellrean.id/belajar-auth/internal/config"
	"shellrean.id/belajar-auth/internal/middleware"
	"shellrean.id/belajar-auth/internal/repository"
	"shellrean.id/belajar-auth/internal/service"
	"shellrean.id/belajar-auth/internal/sse"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	cacheConnection := repository.NewRedisClient(cnf)

	hub := &dto.Hub{
		NotificationChannel: map[int64]chan dto.NotificationData{},
	}

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)
	templateRepository := repository.NewTemplate(dbConnection)
	topUpRepository := repository.NewTopUp(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	notificationService := service.NewNotification(notificationRepository, templateRepository, hub)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationService)
	midtransService := service.NewMidtrans(cnf)
	topUpService := service.NewTopUp(notificationService, midtransService, topUpRepository, accountRepository, transactionRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, authMid, transactionService)
	api.NewNotification(app, authMid, notificationService)
	api.NewTopUp(app, authMid, topUpService)
	api.NewMidtrans(app, midtransService, topUpService)

	sse.NewNotification(app, authMid, hub)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}

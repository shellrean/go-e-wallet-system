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

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	notificationService := service.NewNotification(notificationRepository, templateRepository, hub)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationService)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, authMid, transactionService)
	api.NewNotification(app, authMid, notificationService)

	sse.NewNotification(app, authMid, hub)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}

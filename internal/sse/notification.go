package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"shellrean.id/belajar-auth/dto"
)

type notificationSse struct {
	hub *dto.Hub
}

func NewNotification(app *fiber.App, authMid fiber.Handler, hub *dto.Hub) {
	h := notificationSse{
		hub: hub,
	}

	app.Get("sse/notification-stream", authMid, h.StreamNotification)
}

func (n notificationSse) StreamNotification(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")

	user, ok := ctx.Locals("x-user").(dto.UserData)
	if !ok {
		return ctx.SendStatus(401)
	}
	ch := n.hub.Register(user.ID)

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer n.hub.Unregister(user.ID, ch)

		event := fmt.Sprintf("event: %s\n"+
			"data: \n\n", "initial")
		if _, err := fmt.Fprint(w, event); err != nil {
			return
		}
		if err := w.Flush(); err != nil {
			return
		}

		for notification := range ch {
			data, err := json.Marshal(notification)
			if err != nil {
				continue
			}

			event = fmt.Sprintf("event: %s\n"+
				"data: %s\n\n", "notification-updated", data)

			if _, err := fmt.Fprint(w, event); err != nil {
				return
			}
			if err := w.Flush(); err != nil {
				return
			}
		}
	})

	return nil
}

package router

import (
	"advent-calendar/internal/handler"
	"advent-calendar/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(r fiber.Router) {
	users := r.Group("/users")
	users.Post("/login", handler.Login)
	users.Patch("/confirm", handler.ConfirmRegister)
	users.Get("/check", middleware.AuthMiddleware, handler.Check)
	users.Patch("/refresh", handler.Refresh)
	users.Post("/subscribe", handler.Subscribe)
	users.Delete("/subscribe", handler.UnSubscribe)
	users.Get("/subscribe", handler.UnSubscribe)
}

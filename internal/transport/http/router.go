package http

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *ActHandler) {
	// health
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})
	// acts
	app.Post("/acts", h.Create)
	app.Get("/acts", h.List)
	app.Get("/acts/:id", h.GetByID)
	app.Delete("/acts/:id", h.Delete)

	// metrics
	app.Get("/metrics", MetricsHandler)
}

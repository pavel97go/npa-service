package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/pavel97go/npa-service/internal/repo/memory"
	thttp "github.com/pavel97go/npa-service/internal/transport/http"
	"github.com/pavel97go/npa-service/internal/usecase"
)

func BuildServer() *fiber.App {
	// repo (in-memory)
	r := memory.NewActRepo()
	// usecase
	uc := usecase.NewActUsecase(r)
	// handlers
	h := thttp.NewActHandler(uc)

	// fiber app + middlewares
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(thttp.MetricsMiddleware)

	// routes
	thttp.RegisterRoutes(app, h)

	return app
}

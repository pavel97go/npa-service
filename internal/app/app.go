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
	r := memory.NewActRepo()
	uc := usecase.NewActUsecase(r)
	h := thttp.NewActHandler(uc)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(thttp.MetricsMiddleware)

	thttp.RegisterRoutes(app, h)

	return app
}

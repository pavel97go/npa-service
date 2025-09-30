package http

import (
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

var totalRequests int64

func MetricsMiddleware(c *fiber.Ctx) error {
	atomic.AddInt64(&totalRequests, 1)
	return c.Next()
}

func MetricsHandler(c *fiber.Ctx) error {
	v := atomic.LoadInt64(&totalRequests)
	return c.JSON(fiber.Map{"requests_total": v})
}

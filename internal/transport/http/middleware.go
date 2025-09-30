package http

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	totalRequests int64
	muMetrics     sync.Mutex
)

func incRequests() {
	muMetrics.Lock()
	totalRequests++
	muMetrics.Unlock()
}

func readRequests() int64 {
	muMetrics.Lock()
	v := totalRequests
	muMetrics.Unlock()
	return v
}

func MetricsMiddleware(c *fiber.Ctx) error {
	incRequests()
	return c.Next()
}

func MetricsHandler(c *fiber.Ctx) error {
	v := readRequests()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"requests_total": v})
}

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Act struct {
	ID    int
	Title string
	Type  string
	Date  time.Time
}
type CreateActReq struct {
	Title string `json:"title"`
	Type  string `json:"type`
	Date  string `json:"date"`
}

var errBadDate = errors.New("bad date format, want YYYY-MM-DD")
var acts []Act
var lasID int
var mu sync.RWMutex
var allowedTypes = map[string]struct{}{
	"constitution": {},
	"federal_law":  {},
	"decree":       {},
	"order":        {},
	"other":        {},
}

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errBadDate
	}
	return t, nil
}

func main() {
	fmt.Println(parseDate("2025-10-01")) // должна вернуть дату и <nil>
	fmt.Println(parseDate("01-10-2025")) // должна вернуть zero time и errBadDate
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}

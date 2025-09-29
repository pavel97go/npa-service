package main

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

//
// ===== Domain models =====
//

type Act struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Type  string    `json:"type"`
	Date  time.Time `json:"date"`
}

type CreateActReq struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	Date  string `json:"date"`
}

//
// ===== Errors & in-memory state =====
//

var (
	// public validation errors
	ErrBadID    = errors.New("bad ID")
	errBadTitle = errors.New("title is required")
	errBadType  = errors.New("bad act type")
	errBadDate  = errors.New("bad date format, want YYYY-MM-DD")

	// domain error
	errNotFound = errors.New("act not found")

	// in-memory storage (protected by mu)
	acts   []Act
	lastID int
	mu     sync.RWMutex

	// whitelist for Type
	allowedTypes = map[string]struct{}{
		"constitution": {},
		"federal_law":  {},
		"decree":       {},
		"order":        {},
		"other":        {},
	}
)

//
// ===== Helpers / Utilities (pure functions) =====
//

// case-insensitive substring check (used in listing filters later)
func ContainsFold(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// parse date from yyyy-mm-dd
func parseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errBadDate
	}
	return t, nil
}

// input validation for CreateActReq
func validateInput(in CreateActReq) error {
	if strings.TrimSpace(in.Title) == "" {
		return errBadTitle
	}
	if _, ok := allowedTypes[in.Type]; !ok {
		return errBadType
	}
	if _, err := parseDate(in.Date); err != nil {
		return err
	}
	return nil
}

// SafeAtoi enforces business rule: id must be > 0
func SafeAtoi(s string) (int, error) {
	if strings.TrimSpace(s) == "" {
		return 0, ErrBadID
	}
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		return 0, ErrBadID
	}
	return id, nil
}

//
// ===== Core (business logic, works with in-memory state) =====
//

// createAct: validation + assign ID + append to store
func createAct(in CreateActReq) (Act, error) {
	if err := validateInput(in); err != nil {
		return Act{}, err
	}
	d, _ := parseDate(in.Date)

	mu.Lock()
	defer mu.Unlock()

	lastID++
	act := Act{
		ID:    lastID,
		Title: strings.TrimSpace(in.Title),
		Type:  in.Type,
		Date:  d,
	}
	acts = append(acts, act)
	return act, nil
}

// GetActByID: read path, RLock (no mutation), linear search
func GetActByID(id int) (Act, error) {
	mu.RLock()
	defer mu.RUnlock()
	for _, a := range acts {
		if a.ID == id {
			return a, nil
		}
	}
	return Act{}, errNotFound
}

func listActs(filterType, q string) []Act {
	mu.RLock()
	defer mu.RUnlock()

	out := []Act{}
	for _, act := range acts {
		if filterType != "" && act.Type != filterType {
			continue
		}
		if q != "" && !ContainsFold(act.Title, q) {
			continue
		}
		out = append(out, act)
	}
	return out
}
func deleteAct(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for i := range acts {
		if acts[i].ID == id {
			acts = append(acts[:i], acts[i+1:]...)
			return nil
		}
	}
	return errNotFound
}

//
// ===== Transport / HTTP (Fiber handlers, error mapping) =====
//

func main() {
	app := fiber.New()

	// health
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// POST /acts — create Act (maps validation errors to 400)
	app.Post("/acts", func(c *fiber.Ctx) error {
		var in CreateActReq
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}
		act, err := createAct(in)
		if err != nil {
			switch err {
			case errBadTitle, errBadDate, errBadType:
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
			}
		}
		return c.Status(fiber.StatusCreated).JSON(act)
	})

	// GET /acts/:id — find by id (400 invalid id, 404 not found)
	app.Get("/acts/:id", func(c *fiber.Ctx) error {
		s := c.Params("id")
		id, err := SafeAtoi(s)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}
		act, err := GetActByID(id)
		if err == errNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
		}
		return c.Status(fiber.StatusOK).JSON(act)
	})

	app.Get("/acts", func(c *fiber.Ctx) error {
		ftype := c.Query("type")
		q := c.Query("q")
		items := listActs(ftype, q)
		return c.JSON(fiber.Map{
			"total": len(items),
			"items": items,
		})

	})
	app.Delete("/acts/:id", func(c *fiber.Ctx) error {
		s := c.Params("id")

		id, err := SafeAtoi(s)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		if err := deleteAct(id); err != nil {
			if err == errNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}

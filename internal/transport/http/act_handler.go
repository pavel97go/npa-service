package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/pavel97go/npa-service/internal/domain"
	"github.com/pavel97go/npa-service/internal/usecase"
)

type ActHandler struct {
	uc *usecase.ActUsecase
}

func NewActHandler(uc *usecase.ActUsecase) *ActHandler {
	return &ActHandler{uc: uc}
}

// POST /acts
func (h *ActHandler) Create(c *fiber.Ctx) error {
	var in domain.CreateActReq
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	act, err := h.uc.Create(in)
	if err != nil {
		switch err {
		case domain.ErrBadTitle, domain.ErrBadDate, domain.ErrBadType:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(act)
}

// GET /acts/:id
func (h *ActHandler) GetByID(c *fiber.Ctx) error {
	s := c.Params("id")
	id, err := domain.SafeAtoi(s)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
	}

	act, err := h.uc.Get(id)
	if err == domain.ErrNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
	}

	return c.Status(fiber.StatusOK).JSON(act)
}

// GET /acts
func (h *ActHandler) List(c *fiber.Ctx) error {
	ftype := c.Query("type")
	q := c.Query("q")
	items, _ := h.uc.List(ftype, q)

	return c.JSON(fiber.Map{
		"total": len(items),
		"items": items,
	})
}

// DELETE /acts/:id
func (h *ActHandler) Delete(c *fiber.Ctx) error {
	s := c.Params("id")
	id, err := domain.SafeAtoi(s)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
	}

	if err := h.uc.Delete(id); err != nil {
		if err == domain.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

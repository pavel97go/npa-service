package domain

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Act struct {
	ID    int
	Title string
	Type  string
	Date  time.Time
}

type CreateActReq struct {
	Title string
	Type  string
	Date  string
}

var (
	// публичные ошибки валидации/домена
	ErrBadID    = errors.New("bad ID")
	ErrBadTitle = errors.New("title is required")
	ErrBadType  = errors.New("bad act type")
	ErrBadDate  = errors.New("bad date format, want YYYY-MM-DD")
	ErrNotFound = errors.New("act not found")

	// whitelist для Type
	AllowedTypes = map[string]struct{}{
		"constitution": {},
		"federal_law":  {},
		"decree":       {},
		"order":        {},
		"other":        {},
	}
)

// case-insensitive contains
func ContainsFold(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// yyyy-mm-dd -> time.Time
func ParseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, ErrBadDate
	}
	return t, nil
}

// валидация входа для CreateActReq
func ValidateCreateInput(in CreateActReq) error {
	if strings.TrimSpace(in.Title) == "" {
		return ErrBadTitle
	}
	if _, ok := AllowedTypes[in.Type]; !ok {
		return ErrBadType
	}
	if _, err := ParseDate(in.Date); err != nil {
		return err
	}
	return nil
}

// бизнес-правило: id > 0
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

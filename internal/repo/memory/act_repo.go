package memory

import (
	"sync"

	"github.com/pavel97go/npa-service/internal/domain"
)

type ActRepo struct {
	mu     sync.RWMutex
	acts   []domain.Act
	lastID int
}

func NewActRepo() *ActRepo {
	return &ActRepo{}
}

// Create: присвоение ID и сохранение
func (r *ActRepo) Create(a domain.Act) (domain.Act, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	a.ID = r.lastID
	r.acts = append(r.acts, a)
	return a, nil
}

// ByID: поиск по ID
func (r *ActRepo) ByID(id int) (domain.Act, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.acts {
		if a.ID == id {
			return a, nil
		}
	}
	return domain.Act{}, domain.ErrNotFound
}

// List: фильтрация по type и q (title contains, case-insensitive)
func (r *ActRepo) List(filterType, q string) ([]domain.Act, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Act, 0, len(r.acts))
	for _, a := range r.acts {
		if filterType != "" && a.Type != filterType {
			continue
		}
		if q != "" && !domain.ContainsFold(a.Title, q) {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

// Delete: удаление по ID
func (r *ActRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.acts {
		if r.acts[i].ID == id {
			r.acts = append(r.acts[:i], r.acts[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

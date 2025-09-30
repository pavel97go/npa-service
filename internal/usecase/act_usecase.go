package usecase

import (
	"strings"

	"github.com/pavel97go/npa-service/internal/domain"
)

// интерфейс репозитория
type ActRepository interface {
	Create(a domain.Act) (domain.Act, error)
	ByID(id int) (domain.Act, error)
	List(filterType, q string) ([]domain.Act, error)
	Delete(id int) error
}

type ActUsecase struct {
	repo ActRepository
}

func NewActUsecase(r ActRepository) *ActUsecase {
	return &ActUsecase{repo: r}
}

// повторяет твою логику: валидация -> парс даты -> TrimSpace -> repo.Create
func (uc *ActUsecase) Create(in domain.CreateActReq) (domain.Act, error) {
	if err := domain.ValidateCreateInput(in); err != nil {
		return domain.Act{}, err
	}
	d, _ := domain.ParseDate(in.Date)
	a := domain.Act{
		Title: strings.TrimSpace(in.Title),
		Type:  in.Type,
		Date:  d,
	}
	return uc.repo.Create(a)
}

func (uc *ActUsecase) Get(id int) (domain.Act, error) {
	if id <= 0 {
		return domain.Act{}, domain.ErrBadID
	}
	return uc.repo.ByID(id)
}

func (uc *ActUsecase) List(filterType, q string) ([]domain.Act, error) {
	return uc.repo.List(filterType, q)
}

func (uc *ActUsecase) Delete(id int) error {
	if id <= 0 {
		return domain.ErrBadID
	}
	return uc.repo.Delete(id)
}

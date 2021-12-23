package usecases

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain"
)

type ExpressionInteractor struct {
	ExpressionRepo domain.ExpressionRepo
}

func NewExpressionInteractor(er domain.ExpressionRepo) *ExpressionInteractor {
	return &ExpressionInteractor{er}
}

func (i *ExpressionInteractor) Search(langCode string, value string) ([]*domain.Expression, error) {
	expressions, err := i.ExpressionRepo.Search(langCode, value)
	if err != nil {
		return nil, err
	}

	return expressions, nil
}

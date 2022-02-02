package usecases

import (
	"strings"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
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

func (i *ExpressionInteractor) GetTranscriptionParts(expressionId *valueobject.ID) ([]*domain.TranscriptionPart, error) {
	expression, err := i.ExpressionRepo.Get(expressionId)
	if err != nil {
		return nil, err
	}

	exprParts := strings.Split(expression.Value, "")

	parts, err := i.ExpressionRepo.GetTranscriptionParts(expressionId, exprParts)
	if err != nil {
		return nil, err
	}

	return parts, nil
}

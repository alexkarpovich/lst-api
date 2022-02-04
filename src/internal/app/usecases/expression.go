package usecases

import (
	"fmt"
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

func (i *ExpressionInteractor) GetTranscriptionMap(expressionId *valueobject.ID, typeId *valueobject.ID) (map[string][]*domain.TranscriptionItem, error) {
	expression, err := i.ExpressionRepo.Get(expressionId)
	if err != nil {
		return nil, err
	}

	exprParts := strings.Split(expression.Value, "")

	transcriptionMap, err := i.ExpressionRepo.GetTranscriptionMap(typeId, exprParts)
	if err != nil {
		return nil, err
	}

	fmt.Print(transcriptionMap)

	// transcriptionParts := []*domain.TranscriptionPart{}

	// for _, exprValue := range exprParts {
	// 	part := &domain.TranscriptionPart{

	// 	}
	// }

	return transcriptionMap, nil
}

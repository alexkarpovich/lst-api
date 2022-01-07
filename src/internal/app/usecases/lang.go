package usecases

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain"
)

type LangInteractor struct {
	LangRepo domain.LangRepo
}

func NewLangInteractor(lr domain.LangRepo) *LangInteractor {
	return &LangInteractor{lr}
}

func (i *LangInteractor) List() ([]*domain.Language, error) {
	languages, err := i.LangRepo.List()
	if err != nil {
		return nil, err
	}

	return languages, nil
}

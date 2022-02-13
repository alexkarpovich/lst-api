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

func (i *LangInteractor) ListTranscriptionTypes(langCode string) ([]*domain.TranscriptionType, error) {
	transcriptionTypes, err := i.LangRepo.ListTranscriptionTypes(langCode)
	if err != nil {
		return nil, err
	}

	return transcriptionTypes, nil
}

package usecases

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TranslationInteractor struct {
	TranslationRepo domain.TranslationRepo
}

func NewTranslationInteractor(tr domain.TranslationRepo) *TranslationInteractor {
	return &TranslationInteractor{tr}
}

func (i *TranslationInteractor) AttachTranscription(actorId *valueobject.ID, translationId *valueobject.ID, transcriptionId *valueobject.ID) error {
	err := i.TranslationRepo.AttachTranscription(translationId, transcriptionId)
	if err != nil {
		return err
	}

	return nil
}

func (i *TranslationInteractor) DetachTranscription(actorId *valueobject.ID, translationId *valueobject.ID, transcriptionId *valueobject.ID) error {
	err := i.TranslationRepo.DetachTranscription(translationId, transcriptionId)
	if err != nil {
		return err
	}

	return nil
}

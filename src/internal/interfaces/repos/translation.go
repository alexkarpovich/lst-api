package repos

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type TranslationRepo struct {
	db db.DB
}

func NewTranslationRepo(db db.DB) *TranslationRepo {
	return &TranslationRepo{db}
}

func (r *TranslationRepo) AttachTranscription(translationId *valueobject.ID, transcriptionId *valueobject.ID) error {
	query := `
		INSERT INTO translation_transcription (translation_id, transcription_id)
		VALUES ($1, $2)
	`
	_, err := r.db.Db().Exec(query, translationId, transcriptionId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TranslationRepo) DetachTranscription(translationId *valueobject.ID, transcriptionId *valueobject.ID) error {
	query := `
		DELETE FROM translation_transcription
		WHERE translation_id=$1 AND transcription_id=$2
	`
	_, err := r.db.Db().Exec(query, translationId, transcriptionId)
	if err != nil {
		return err
	}

	return nil
}

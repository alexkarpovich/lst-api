package domain

import "github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

type Translation struct {
	Id             *valueobject.ID `json:"id" db:"id"`
	Target         *Expression     `json:"target"`
	Native         *Expression     `json:"native"`
	Transcriptions []string        `json:"transcriptions" db:"transcriptions"`
	Comment        string          `json:"comment" db:"comment"`
}

type TranslationRepo interface {
	AttachTranscription(*valueobject.ID, *valueobject.ID) error
	DetachTranscription(*valueobject.ID, *valueobject.ID) error
}

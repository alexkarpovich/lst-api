package domain

import "github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

type Expression struct {
	Id             *valueobject.ID      `json:"id" db:"id"`
	LangCode       string               `json:"langCode" db:"lang"`
	Value          string               `json:"value" db:"value"`
	Transcriptions []*TranscriptionItem `json:"transcriptions"`
	Lang           *Language            `json:"lang"`
}

type Transcription struct {
	Id    *valueobject.ID `json:"id" db:"id"`
	Type  *valueobject.ID `json:"type" db:"type"`
	Value string          `json:"value" db:"value"`
}

type TranscriptionItem struct {
	Id    *valueobject.ID `json:"id" db:"id"`
	Value string          `json:"value" db:"value"`
}

type TranscriptionPart struct {
	Expression     Expression           `json:"expression"`
	Transcriptions []*TranscriptionItem `json:"transcriptions"`
}

type ExpressionRepo interface {
	Create(*Expression) (*Expression, error)
	Get(*valueobject.ID) (*Expression, error)
	Search(string, string) ([]*Expression, error)
	CreateTranscription(*valueobject.ID, Transcription) (*Transcription, error)
	GetTranscriptionMap(*valueobject.ID, *valueobject.ID) (map[string][]*TranscriptionItem, error)
}

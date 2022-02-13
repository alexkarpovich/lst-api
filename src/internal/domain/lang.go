package domain

import "github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

type Language struct {
	Code       string `json:"code" db:"code"`
	IsoName    string `json:"isoName" db:"iso_name"`
	NativeName string `json:"nativeName" db:"native_name"`
}

type TranscriptionType struct {
	Id   *valueobject.ID `json:"id" db:"id"`
	Name string          `json:"name" db:"name"`
}

type LangRepo interface {
	List() ([]*Language, error)
	ListTranscriptionTypes(string) ([]*TranscriptionType, error)
}

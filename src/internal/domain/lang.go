package domain

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type Language struct {
	Id         *valueobject.ID `json:"id"`
	Code       string          `json:"code"`
	IsoName    string          `json:"isoName"`
	NativeName string          `json:"nativeName"`
}

type LangRepo interface {
	List() []*Language
}

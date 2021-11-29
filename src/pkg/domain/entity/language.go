package entity

import (
	"pkg/domain/valueobject"
)

type Language struct {
	Id         valueobject.ID `json:"id"`
	Code       string         `json:"code"`
	IsoName    string         `json:"isoName"`
	NativeName string         `json:"nativeName"`
}

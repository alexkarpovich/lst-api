package app

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type GroupStatus uint

const (
	GroupActive GroupStatus = iota
	GroupDeleted
)

type Group struct {
	Id             *valueobject.ID `json:"id" db:"id"`
	TargetLangCode string          `json:"targetLangCode" db:"target_lang"`
	NativeLangCode string          `json:"nativeLangCode" db:"native_lang"`
	Name           string          `json:"name" db:"name"`
	Status         GroupStatus     `json:"status" db:"status"`
}

type GroupRepo interface {
	Create(*Group) (*Group, error)
	List(*valueobject.ID) ([]*Group, error)
	MarkAsDeleted(*valueobject.ID) error
	AttachUser(*valueobject.ID, *valueobject.ID, UserRole) error
	DetachUser(*valueobject.ID, *valueobject.ID) error
}

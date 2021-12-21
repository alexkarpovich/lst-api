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
	Id           *valueobject.ID `json:"id" db:"id"`
	TargetLangId *valueobject.ID `json:"targetLangId" db:"target_lang_id"`
	NativeLangId *valueobject.ID `json:"nativeLangId" db:"native_lang_id"`
	Name         string          `json:"name" db:"name"`
	Status       GroupStatus     `json:"status" db:"status"`
}

type GroupRepo interface {
	Create(*Group) (*Group, error)
	List(*valueobject.ID) ([]*Group, error)
	MarkAsDeleted(*valueobject.ID) error
	AttachUser(*valueobject.ID, *valueobject.ID, UserRole) error
	DetachUser(*valueobject.ID, *valueobject.ID) error
}

package app

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type GroupStatus uint

const (
	GroupActive GroupStatus = iota
	GroupDeleted
)

type GroupMember struct {
	Id       *valueobject.ID `json:"id" db:"id"`
	Username string          `json:"username" db:"username"`
	Role     UserRole        `json:"role" db:"role"`
}

type Group struct {
	Id             *valueobject.ID `json:"id" db:"id"`
	TargetLangCode string          `json:"targetLangCode" db:"target_lang"`
	NativeLangCode string          `json:"nativeLangCode" db:"native_lang"`
	Name           string          `json:"name" db:"name"`
	Status         GroupStatus     `json:"status" db:"status"`
	Members        []*GroupMember  `json:"members"`
}

type GroupRepo interface {
	Create(*Group) (*Group, error)
	List(*valueobject.ID) ([]*Group, error)
	MarkAsDeleted(*valueobject.ID) error
	AttachMember(*valueobject.ID, *valueobject.ID, UserRole) error
	DetachMember(*valueobject.ID, *valueobject.ID) error
}

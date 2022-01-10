package app

import (
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type GroupStatus uint

const (
	GroupActive GroupStatus = iota
	GroupDeleted
)

type GroupMember struct {
	Id             *valueobject.ID `json:"id" db:"id"`
	Username       string          `json:"username" db:"username"`
	Role           UserRole        `json:"role" db:"role"`
	Status         MemberStatus    `json:"status" db:"status"`
	Token          string          `db:"token"`
	TokenExpiresAt time.Time       `db:"token_expires_at"`
}

type Group struct {
	Id             *valueobject.ID `json:"id" db:"id"`
	TargetLangCode string          `json:"targetLangCode" db:"target_lang"`
	NativeLangCode string          `json:"nativeLangCode" db:"native_lang"`
	Name           string          `json:"name" db:"name"`
	Status         GroupStatus     `json:"status" db:"status"`
	IsUntouched    bool            `json:"isUntouched"`
	Members        []*GroupMember  `json:"members"`
}

type GroupRepo interface {
	Create(*valueobject.ID, *Group) (*Group, error)
	Update(*Group) error
	List(*valueobject.ID) ([]*Group, error)
	MarkAsDeleted(*valueobject.ID) error
	FindMemberById(*valueobject.ID, *valueobject.ID) (*GroupMember, error)
	FindMemberByToken(string) (*valueobject.ID, *GroupMember, error)
	AttachUser(*valueobject.ID, GroupMember) error
	DetachMember(*valueobject.ID, *valueobject.ID) error
	UpdateMember(*valueobject.ID, GroupMember) error
}

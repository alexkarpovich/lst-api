package app

import (
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

	"golang.org/x/crypto/bcrypt"
)

// UserRole ...
type UserRole uint

const (
	UserAdmin UserRole = iota
	UserReader
	UserEditor
)

// UserStatus is used to identiry user statuses
type UserStatus uint

const (
	UserInactive UserStatus = iota
	UserUnconfirmed
	UserActive
	UserDeleted
)

type MemberStatus uint

const (
	MemberPending MemberStatus = iota
	MemberActive
	MemberDeleted
)

type User struct {
	Id                *valueobject.ID          `json:"id" db:"id"`
	Email             valueobject.EmailAddress `json:"email" db:"email"`
	Username          string                   `json:"username" db:"username"`
	EncryptedPassword string                   `json:"-" db:"encrypted_password"`
	FirstName         string                   `json:"firstName" db:"first_name"`
	LastName          string                   `json:"lastName" db:"last_name"`
	Token             string                   `json:"token" db:"token"`
	Role              UserRole                 `json:"role" db:"role"`
	Status            UserStatus               `json:"status" db:"status"`
	TokenExpiresAt    time.Time                `json:"-" db:"token_expires_at"`
	CreatedAt         time.Time                `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time                `json:"updatedAt" db:"updated_at"`
}

type UserRepo interface {
	Create(User) (*User, error)
	Get(*valueobject.ID) (*User, error)
	FindByUsername(string) (*User, error)
	FindByEmail(valueobject.EmailAddress) (*User, error)
	FindByToken(string) (*User, error)
	Delete(*valueobject.ID) error
	Update(User) error
}

func (user *User) SetPassword(plainPassword string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 14)
	user.EncryptedPassword = string(bytes)
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))

	return err == nil
}

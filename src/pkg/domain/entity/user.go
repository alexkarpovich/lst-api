package entity

import (
	"pkg/domain/valueobject"

	"golang.org/x/crypto/bcrypt"
)

// UserStatus is used to identiry user statuses
type UserStatus uint

const (
	UserInactive UserStatus = iota
	UserActive
	UserDeleted
)

type User struct {
	Id                valueobject.ID           `json:"id"`
	Email             valueobject.EmailAddress `json:"email"`
	EncryptedPassword string                   `json:"-"`
	FirstName         string                   `json:"firstName"`
	LastName          string                   `json:"lastName"`
	Status            UserStatus               `json:"status"`
}

func (user *User) SetPassword(plainPassword string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 14)
	user.EncryptedPassword = string(bytes)
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))

	return err == nil
}

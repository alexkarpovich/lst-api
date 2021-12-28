package services

import "github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

type EmailService interface {
	SendSignup(valueobject.EmailAddress, string) error
}

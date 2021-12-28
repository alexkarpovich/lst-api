package services

import (
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/services/email"
)

type Services struct {
	Email services.EmailService
}

func NewServices() *Services {
	return &Services{
		Email: &email.EmailService{},
	}
}

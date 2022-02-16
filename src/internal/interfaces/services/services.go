package services

import (
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/repos"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/services/email"
)

type Services struct {
	Email    services.EmailService
	Training services.TrainingService
}

func NewServices(repos *repos.Repos) *Services {
	return &Services{
		Email: &email.EmailService{},
	}
}

package usecases

import (
	"errors"
	"log"
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/pkg"
)

const tokenLength = 32

type Registrant struct {
	Username string
	Email    string
	Password string
}

type AuthInteractor struct {
	UserRepo     app.UserRepo
	EmailService services.EmailService
}

func NewAuthInteractor(ur app.UserRepo, es services.EmailService) *AuthInteractor {
	return &AuthInteractor{ur, es}
}

func (i *AuthInteractor) CreateRegistrant(r *Registrant) (*app.User, error) {
	registrant := &app.User{
		Email:          valueobject.EmailAddress(r.Email),
		Username:       r.Username,
		Token:          pkg.RandomString(tokenLength),
		Status:         app.UserUnconfirmed,
		TokenExpiresAt: time.Now().Add(12 * time.Hour).UTC(),
	}
	registrant.SetPassword(r.Password)

	registrant, err := i.UserRepo.Create(registrant)
	if err != nil {
		return nil, err
	}

	go i.EmailService.SendSignup(registrant.Email, registrant.Token)

	return registrant, nil
}

func (i *AuthInteractor) ConfirmEmail(token string) error {
	registrant, err := i.UserRepo.FindByToken(token)
	if err != nil {
		log.Println(err)
		return err
	}

	registrant.Status = app.UserActive

	if err := i.UserRepo.Update(registrant); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (i *AuthInteractor) FindUserByEmailPassword(email string, password string) (*app.User, error) {
	user, err := i.UserRepo.FindByEmail(valueobject.EmailAddress(email))
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

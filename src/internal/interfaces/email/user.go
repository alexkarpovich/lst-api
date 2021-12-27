package email

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

func SendSignup(email valueobject.EmailAddress, token string) {
	subject := "Подтверждение регистрации"
	from := "alexsure.k@gmail.com"

	data := make(map[string]interface{})
	data["Token"] = token

	SendWithView(
		subject,
		from,
		[]string{string(email)},
		[]string{
			"./assets/email/templates/layout/base.html",
			"./assets/email/templates/auth/signup.html",
		},
		"layout",
		data,
	)
}

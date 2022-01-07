package email

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

func (s *EmailService) SendGroupInvitation(email valueobject.EmailAddress, token string) error {
	subject := "Group Invitation"
	from := "admin@akarpovich.com"

	data := make(map[string]interface{})
	data["Token"] = token

	return s.SendWithView(
		subject,
		from,
		[]string{string(email)},
		[]string{
			"./assets/email/layout/base.html",
			"./assets/email/group/invitation.html",
		},
		"layout",
		data,
	)
}

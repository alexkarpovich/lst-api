package email

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct{}

func (s *EmailService) Send(subject string, from string, recipients []string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"))

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (s *EmailService) SendWithView(subject string, from string, recipients []string, views []string, layout string, data interface{}) error {
	t, err := template.ParseFiles(views...)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, layout, data); err != nil {
		log.Fatal(err)
		return err
	}

	return s.Send(subject, from, recipients, tpl.String())
}

package email

import (
	"bytes"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func Send(subject string, from string, recipients []string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMT_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendWithView(subject string, from string, recipients []string, views []string, layout string, data interface{}) {
	t, err := template.ParseFiles(views...)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, layout, data); err != nil {
		panic(err)
	}

	Send(subject, from, recipients, tpl.String())
}

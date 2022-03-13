package main

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"strconv"

	gomail "gopkg.in/gomail.v2"
)

func SendMail(subject string, templateHtml string, rcpt []string, data interface{}) {

	if subject == "" {
		subject = os.Getenv("MAIL_SUBJECT")
	}

	if templateHtml == "" {
		templateHtml = "default.html"
	}

	var err error

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	t := template.New(templateHtml)

	templateHtmlFile := wd + "/html/" + templateHtml
	t, err = t.ParseFiles(templateHtmlFile)
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_SENDER"))
	m.SetHeader("To", rcpt...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", result)
	// m.Attach(templateHtml) // attach whatever you want

	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), mailPort, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASS"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

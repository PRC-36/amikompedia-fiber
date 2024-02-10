package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,

	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{name: name, fromEmailAddress: fromEmailAddress, fromEmailPassword: fromEmailPassword}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,

) error {
	e := email.NewEmail()
	e.From = sender.name + " <" + sender.fromEmailAddress + ">"
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, attachFile := range attachFiles {
		_, err := e.AttachFile(attachFile)
		if err != nil {
			return fmt.Errorf("error when attach file: %w", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	err := e.Send(smtpServerAddress, smtpAuth)

	if err != nil {
		return fmt.Errorf("error when send email: %w", err)
	}

	return nil
}

package mail

//
//import (
//	"fmt"
//	"gopkg.in/mail.v2"
//)
//
//type EmailSender interface {
//	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
//}
//
//type TitanSender struct {
//	name              string
//	fromEmailAddress  string
//	fromEmailPassword string
//}
//
//func NewTitanSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
//	return &TitanSender{
//		name:              name,
//		fromEmailAddress:  fromEmailAddress,
//		fromEmailPassword: fromEmailPassword,
//	}
//}
//
//func (sender *TitanSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
//	m := mail.NewMessage()
//	m.SetHeader("From", sender.fromEmailAddress)
//	m.SetHeader("To", to...)
//	m.SetHeader("Cc", cc...)
//	m.SetHeader("Bcc", bcc...)
//	m.SetHeader("Subject", subject)
//	m.SetBody("text/html", content)
//
//	for _, attachFile := range attachFiles {
//		m.Attach(attachFile)
//	}
//
//	d := mail.NewDialer("smtp.titan.email", 587, sender.fromEmailAddress, sender.fromEmailPassword)
//
//	// Send the email
//	if err := d.DialAndSend(m); err != nil {
//		return fmt.Errorf("error when sending email: %w", err)
//	}
//
//	return nil
//}

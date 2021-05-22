package mail

import (
	"crypto/tls"
	"strconv"

	"github.com/apex/log"
	"github.com/inheaden/simple-mail-api/pkg/config"
	gomail "gopkg.in/mail.v2"
)

func SendMailWithoutFrom(to string, emailSubject string, emailBody string) error {

	m := gomail.NewMessage()

	mailConfig := config.GetMailConfig()

	m.SetHeader("From", mailConfig.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/plain", emailBody)

	err := Sendmail(m)
	if err != nil {
		return err
	}

	log.Debugf("Email to %s was sent", to)

	return nil
}

func SendMailWithFrom(to string, emailSubject string, emailBody string, from string) error {

	m := gomail.NewMessage()

	mailConfig := config.GetMailConfig()

	m.SetHeader("From", mailConfig.SMTPFrom)
	m.SetHeader("Reply-To", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/plain", emailBody)

	err := Sendmail(m)
	if err != nil {
		return err
	}

	log.Debugf("Email to %s was sent", to)

	return nil
}

// Sendmail sends a single mail
func Sendmail(m *gomail.Message) error {
	mailConfig := config.GetMailConfig()

	port, _ := strconv.Atoi(mailConfig.SMTPPort)
	d := gomail.NewDialer(mailConfig.SMTPURL, port, mailConfig.Username, mailConfig.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: mailConfig.SMTPURL}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

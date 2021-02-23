package mail

import (
	"crypto/tls"
	"log"
	"strconv"

	gomail "gopkg.in/mail.v2"
	"inheaden.io/services/simple-mail-api/pkg/config"
)

// Sendmail sends a single mail
func Sendmail(to string, emailSubject string, emailBody string) error {
	m := gomail.NewMessage()

	mailConfig := config.GetMailConfig()

	m.SetHeader("From", mailConfig.SmtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/plain", emailBody)

	port, _ := strconv.Atoi(mailConfig.SmtpPort)
	d := gomail.NewDialer(mailConfig.SmtpURL, port, mailConfig.Username, mailConfig.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: mailConfig.SmtpURL}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	log.Print("Email send")

	return nil
}

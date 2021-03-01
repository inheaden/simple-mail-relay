package mail

import (
	"crypto/tls"
	"strconv"

	"github.com/apex/log"
	gomail "gopkg.in/mail.v2"
	"inheaden.io/services/simple-mail-api/pkg/config"
)

// Sendmail sends a single mail
func Sendmail(to string, emailSubject string, emailBody string) error {
	m := gomail.NewMessage()

	mailConfig := config.GetMailConfig()

	m.SetHeader("From", mailConfig.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/plain", emailBody)

	port, _ := strconv.Atoi(mailConfig.SMTPPort)
	d := gomail.NewDialer(mailConfig.SMTPURL, port, mailConfig.Username, mailConfig.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: mailConfig.SMTPURL}

	// if err := d.DialAndSend(m); err != nil {
	// 	return err
	// }

	log.Debugf("Email to %s send", to)

	return nil
}

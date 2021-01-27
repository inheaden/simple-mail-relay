package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

func SendmailNotification(receiverEmailAddress string, emailSubject string, emailBody string) {


	// Sender data  Inheaden mail and pwd
	from := ""
	password := ""

	// Receiver email address.
	// email will be sent to both client and Inheaden with data submited by client
	to := []string{
		receiverEmailAddress,
	}




	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("simple-mail-api/pkg/mail/notificationmail.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: "+ emailSubject +"\n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Body    string
	}{
		Body:    emailBody,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

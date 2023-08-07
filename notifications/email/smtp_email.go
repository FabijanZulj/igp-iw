package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

// SimpleSMTPEmailSender sends emails with SMTP using plain auth
// and implements EmailSender interface
type SimpleSMTPEmailSender struct {
	SmtpHost string
	SmtpPort string
	From     string
	Subject  string
	Username string
	Password string
}

// SendMail uses SMTP plain auth, the message is sent with Content-Type text/html
// so it expects a html string in the `msg` argument
func (es *SimpleSMTPEmailSender) SendMail(to, subject string, msg []byte) error {
	log.Println("Sending email")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	toRaw := fmt.Sprintf("To: %s\r\n", to)
	rawSubject := fmt.Sprintf("Subject: %s\r\n", subject)
	body := fmt.Sprintf("\r\n%s\r\n", string(msg))

	fullMsg := []byte(toRaw + rawSubject + mime + body)

	auth := smtp.PlainAuth("", es.Username, es.Password, es.SmtpHost)

	err := smtp.SendMail(es.SmtpHost+":"+es.SmtpPort, auth, es.From, []string{to}, []byte(fullMsg))
	if err != nil {
		log.Println("Error sending email with Plain Auth SMTP " + err.Error())
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}

// CreateMail tries to find a template in /templates folder
// executes the template with the given MailData and returns the parsed and executed template
func (es *SimpleSMTPEmailSender) CreateMail(ed EmailData) []byte {
	t, err := template.ParseFiles("templates/" + ed.Template)
	if err != nil {
		log.Println("Could not parse files " + err.Error())
		return nil
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, ed.MailData); err != nil {
		log.Println("Could not execute template " + err.Error())
		return nil
	}
	return buf.Bytes()
}

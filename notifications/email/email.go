package email

// EmailData is representation of an email that gets compiled
// to email string/byte[] that can be sent
type EmailData struct {
	MailData map[string]any
	To       string
	Template string
}

// EmailSender provides the functionality of email creation/parsing from EmailData
// and sending emails
type EmailSender interface {
	CreateMail(EmailData) []byte
	SendMail(string, string, []byte) error
}

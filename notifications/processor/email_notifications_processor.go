package processor

import "githubn.com/igp-iw-notifications/email"

// EmailNotificationsProcessor implements the NotificationProcessor interface and provides a way of
// processing email notifications using the `email` package.
// The processor contains an implementation if a mail sending service
type EmailNotificationsProcessor struct {
	email email.EmailSender
}

// NewEmailNotificationsProcessor creates a email processor using the given email sending service `es`
func NewEmailNotificationsProcessor(es email.EmailSender) *EmailNotificationsProcessor {
	return &EmailNotificationsProcessor{
		email: es,
	}
}

// Process is the function that processes the email notification by creating the email string
// using the template that has to be passed as additional data
// The Process function does not handle errors while sending the email and just skips if it is unsucessful
func (evp *EmailNotificationsProcessor) Process(nd *NotificationData) {
	template, ok := nd.RawData["template"].(string)
	if !ok {
		return
	}
	subject, ok := nd.RawData["subject"].(string)
	if !ok {
		return
	}

	emailTemplate := evp.email.CreateMail(email.EmailData{
		MailData: nd.RawData,
		To:       nd.Target,
		Template: template,
	})
	evp.email.SendMail(nd.Target, subject, emailTemplate)
}

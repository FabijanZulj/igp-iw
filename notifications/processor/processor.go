package processor

// NotificationData contains all the data needed for processing,
// additional data that is processor implementation specific is passed as RawData
type NotificationData struct {
	// RawData contains additional implementation specific data or config for each processor
	RawData map[string]any `json:"rawData"`
	// NotificationType is the type of the notification that has to be dispatched
	// - (eg. EMAIL_NOTIFICATION, SMS_NOTIFICATION)
	NotificationType string `json:"notificationType"`
	// Initiator is the identifier of the actor who dispatched the notification
	Initiator string `json:"initiator"`
	// Target is the target of the notification (eg. email address, sms number, webhook url)
	Target string `json:"target"`
}

// NotificationProcessor provides an interface for various ways of processing notifications (eg. email, webhook, SMS)
type NotificationProcessor interface {
	Process(*NotificationData)
}

package notifications

import (
	"encoding/json"
	"errors"
	"log"

	"githubn.com/igp-iw-notifications/config"
	"githubn.com/igp-iw-notifications/email"
	"githubn.com/igp-iw-notifications/processor"
)

// RunNotificationHub receives notifications from the provided `notifChan` and dispatches
// them to the approriate processor
func RunNotificationHub(notifChan chan []byte, config *config.Config) {
	for {
		notifData := <-notifChan
		nd := &processor.NotificationData{}
		err := json.Unmarshal(notifData, nd)
		if err != nil {
			log.Println("Error unmarshaling" + err.Error())
			continue
		}
		log.Printf("Successfuly unmarshaled notification with type: %v \n" + nd.NotificationType)
		processor, err := mapTypeToProcessor(nd, config)
		if err != nil {
			log.Printf("Notification not processed:" + err.Error())
			continue
		}
		processor.Process(nd)
	}
}

func mapTypeToProcessor(notificationData *processor.NotificationData, config *config.Config) (processor.NotificationProcessor, error) {
	switch notificationData.NotificationType {
	case "EMAIL_NOTIFICATION":
		return processor.NewEmailNotificationsProcessor(&email.SimpleSMTPEmailSender{
			SmtpHost: config.SmtpHost,
			SmtpPort: config.SmtpPort,
			From:     config.SmtpFrom,
			Username: config.SmtpUsername,
			Password: config.SmtpPassword,
		}), nil
	case "SMS_NOTIFICATION":
		return &processor.SMSNotificationsProcessor{}, nil
	default:
		return nil, errors.New("could not find a processor for given notification type")
	}
}

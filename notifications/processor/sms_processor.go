package processor

import "log"

// SMSNotificationsProcessor is a sample implementation of the NotificationProcessor interface
// that can send SMS notifications
type SMSNotificationsProcessor struct{}

func (snp *SMSNotificationsProcessor) Process(nd *NotificationData) {
	data, ok := nd.RawData["smsData"].(string)
	if !ok {
		log.Println("Sms notification not valid")
		return
	}
	log.Println("================================")
	log.Println("========SAMPLE SMS NOTIF========")
	log.Println(nd.Target)
	log.Println(data)
	log.Println("================================")
}

package main

import (
	notifications "githubn.com/igp-iw-notifications"
	"githubn.com/igp-iw-notifications/config"
	"githubn.com/igp-iw-notifications/consumers"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(1)
	}
	notifChan := make(chan []byte)
	go notifications.RunNotificationHub(notifChan, config)

	consumer := &consumers.KafkaConsumer{}
	metadata := consumers.ConsumerMetadata{
		RawData: map[string]any{
			"brokers":       config.KafkaBrokers,
			"consumerGroup": config.ConsumerGroup,
		},
		Topic: config.KafkaTopic,
	}

	consumer.Consume(metadata, notifChan)
}

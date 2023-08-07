package consumers

// ConsumerMetadata contains implementation specific config fields
// needed for the consumer and the topic to consume from
type ConsumerMetadata struct {
	RawData map[string]any
	Topic   string
}

// NotificationsConsumer defines a consumer that blocks on Consume() call and sends the
// consumed messages to the channel
// (this can be PubSub, Kafka consumer...)
type NotificationsConsumer interface {
	Consume(ConsumerMetadata, chan []byte)
}

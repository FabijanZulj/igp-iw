package notifications

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

// KafkaNotificationPublisher implements the NotificationPublisher interfaces by
// publishing the notifications to a pubsub topic
type KafkaNotificationPublisher struct {
	producer sarama.AsyncProducer
	topic    string
}

// Publish marshals the notification data `nd` and sends it to the topic
// with the key being the Initiator of the notification
func (knp *KafkaNotificationPublisher) Publish(nd NotificationData) error {
	b, err := json.Marshal(nd)
	if err != nil {
		log.Println("Error marshaling notification data" + err.Error())
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: knp.topic,
		Value: sarama.ByteEncoder(b),
		Key:   sarama.StringEncoder(nd.Initiator),
	}
	knp.producer.Input() <- message
	log.Println("Published notification" + string(b))
	return nil
}

// NewKafkaNotificationPublisher creates a new kafka publisher.
// !PANICS! if connecting to any of the brokers is not posible
func NewKafkaNotificationPublisher(kafkaBrokers []string, topic string) (*KafkaNotificationPublisher, error) {
	config := sarama.NewConfig()
	producer := kafkaPublisherRetrying(kafkaBrokers, config)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func(signals chan os.Signal) {
		for {
			select {
			case <-signals:
				log.Println("signal recieved")
				producer.AsyncClose()
				os.Exit(0)
			case x := <-producer.Errors():
				log.Println("Error while publishing" + x.Error())
			}
		}
	}(signals)

	return &KafkaNotificationPublisher{
		producer: producer,
		topic:    topic,
	}, nil
}

func kafkaPublisherRetrying(brokers []string, config *sarama.Config) sarama.AsyncProducer {
	timeout := 10
	retry := 0
	for {
		log.Printf("Trying to connect to kafka brokers : %v \n", brokers)
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err == nil {
			log.Println("Established connection to kafka")
			return producer
		}
		if retry == timeout {
			log.Panicf("Error connecting to kafka after: %v retries", timeout)
		}
		retry++
		time.Sleep(2 * time.Second)
	}
}

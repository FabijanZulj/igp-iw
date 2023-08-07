package consumers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

// KafkaConsumer is an implementation of a notifications consumer that consumes the messages
// from a kafka topic
type KafkaConsumer struct{}

// Consume starts the blocking consuming of notifications from a Kafka topic in a ConsumerGroup
// with a groupID provided in the `cm` ConsumerMetadata argument
// and sends the messages/notifications to a notificationsChannel
func (*KafkaConsumer) Consume(cm ConsumerMetadata, notificationsChannel chan []byte) {
	keepRunning := true
	log.Println("Starting a new Kafka consumer")

	brokers, ok := cm.RawData["brokers"].([]string)
	if !ok {
		log.Panicf("Please provide a list of brokers")
	}
	consumerGroupName, ok := cm.RawData["consumerGroup"].(string)
	if !ok {
		log.Panicf("Please provide the name of consumer group")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer := Consumer{
		ready:                make(chan bool),
		notificationsChannel: notificationsChannel,
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := retryConsumerGroupClient(brokers, consumerGroupName, config)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{cm.Topic}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()
	if err := client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func retryConsumerGroupClient(addrs []string, groupID string, config *sarama.Config) sarama.ConsumerGroup {
	timeout := 10
	retry := 0
	for {
		log.Printf("Trying to connect to kafka broker : %v, consumer group: %v \n", addrs, groupID)
		client, err := sarama.NewConsumerGroup(addrs, groupID, config)
		if err == nil {
			log.Println("Established connection to kafka")
			return client
		}
		if retry == timeout {
			log.Panicf("Error connecting to kafka after: %v retries", timeout)
		}
		retry++
		time.Sleep(2 * time.Second)
	}
}

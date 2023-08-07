package consumers

import (
	"log"

	"github.com/IBM/sarama"
)

// Consumer type implements the ConsumerGroupHandler interface that provides a way to hook in to the
// consumer groupo session lifecycle. It provides a channel that notifies that the setup is done
// and a channel that contains the notifications to further process
type Consumer struct {
	ready                chan bool
	notificationsChannel chan []byte
}

// Setup notifies using the 'ready' channel that the loop is ready to be started
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim starts a consumer loop of ConsumerGroupClaim's Messages().
// It then sends the message to the notifications channel that is further processed by processors
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			consumer.notificationsChannel <- message.Value
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

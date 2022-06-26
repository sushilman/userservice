// Dummy message broker emulator

package messagebroker

import (
	"log"

	"github.com/goccy/go-json"
)

type IMessageBroker interface {
	Publish(string, interface{}) error
}

// This is simply a mock implementation of the message broker emulator
type MessageBroker struct{}

func InitMessageBroker() MessageBroker {
	return MessageBroker{}
}

func (mb MessageBroker) Publish(topic string, message interface{}) error {
	msg, _ := json.Marshal(message)
	log.Printf("Publishing message to: \nTopic: %s \nMessage: %s", topic, string(msg))
	return nil
}

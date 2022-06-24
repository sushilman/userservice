package messagebroker

import (
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
)

type IMessageBroker interface {
	Publish(*zerolog.Logger, string, interface{}) error
}

// This is simply a mock implementation of the message broker emulator
type MessageBroker struct{}

func InitMessageBroker() MessageBroker {
	return MessageBroker{}
}

func (mb MessageBroker) Publish(logger *zerolog.Logger, topic string, message interface{}) error {
	msg, _ := json.Marshal(message)
	logger.Info().Str("topic", topic).Str("message", string(msg)).Msg("Publishing to message broker ")
	return nil
}

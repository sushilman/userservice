package mocks

import "github.com/stretchr/testify/mock"

type MockMessageBroker struct {
	mock.Mock
}

func (m *MockMessageBroker) Publish(topic string, message interface{}) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

package pubsub

import (
	natsgo "github.com/nats-io/nats.go"
)

// IPubSub defines a PubSub does.
type IPubSub interface {
	// Publish sends a message to a topic.
	Publish(topic string, message any) error

	// Subscribe subscribes to a topic and returns a channel for receiving messages.
	Subscribe(topic, queue string, cb func([]byte)) Subscription

	// Unsubscribe unsubscribes from a topic.
	Unsubscribe(topic string) error

	// Close closes the connection to the Pub Sub broker.
	Close() error

	// GetName returns the pubsub name.
	GetName() string

	// GetClient returns the storage client. Use that to interact with the
	// underlying storage client.
	GetClient() any
}

// Subscription is a subscription to a topic.
type Subscription struct {
	// Topic is the subject to subscribe to.
	Topic string

	// Queue is the queue to subscribe to.
	Queue string

	// Callback is the function to call when a message is received.
	Callback func(msg []byte)

	// Channel is the channel to receive messages.
	Channel <-chan []byte
}

// HandleMessage calls the callback function to handle the message.
func (s *Subscription) HandleMessage(msg *natsgo.Msg) {
	s.Callback(msg.Data)
}

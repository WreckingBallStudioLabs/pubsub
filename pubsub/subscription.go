package pubsub

import natsgo "github.com/nats-io/nats.go"

//////
// Const, vars, and types.
//////

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

//////
// Methods.
//////

// HandleMessage calls the callback function to handle the message.
func (s *Subscription) HandleMessage(msg *natsgo.Msg) {
	s.Callback(msg.Data)
}

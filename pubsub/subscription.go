package pubsub

import natsgo "github.com/nats-io/nats.go"

//////
// Const, vars, and types.
//////

// CallBackFunc is the function to call when a message is received.
type CallBackFunc func(msg []byte)

// Subscription is a subscription to a topic.
type Subscription struct {
	// Topic is the subject to subscribe to.
	Topic string

	// Queue is the queue to subscribe to.
	Queue string

	// Callback is the function to call when a message is received.
	Callback CallBackFunc

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

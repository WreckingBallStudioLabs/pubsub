package subscription

import (
	natsgo "github.com/nats-io/nats.go"
	"github.com/thalesfsp/validation"
)

//////
// Const, vars, and types.
//////

// Func is the function to call when a message is received.
type Func func(msg []byte)

// Subscription is a subscription to a topic.
type Subscription struct {
	// Topic is the subject to subscribe to, in the form "v1.meta.created".
	// A "topic" is a way to organize messages.
	Topic string `json:"topic" validate:"required"`

	// Queue is the queue to subscribe to, in the form "v1.meta.created.queue".
	// A "queue" is a way to make sure messages are only delivered to one
	// subscriber at a time.
	Queue string `json:"queue" validate:"required"`

	// Callback is the function to call when a message is received.
	Callback Func `json:"-" validate:"required"`

	// Channel is the channel to receive messages.
	Channel <-chan []byte `json:"-"`
}

//////
// Methods.
//////

// Handle calls the callback function to handle the message.
func (s *Subscription) Handle(msg *natsgo.Msg) {
	s.Callback(msg.Data)
}

//////
// Factory.
//////

// New creates a new subscription. topic and queue should be in the form of the
// following example: "v1.meta.created" and "v1.meta.created.queue".
func New(topic, queue string, callback Func) (*Subscription, error) {
	s := &Subscription{
		Topic:    topic,
		Queue:    queue,
		Callback: callback,
	}

	if err := validation.Validate(s); err != nil {
		return nil, err
	}

	return s, nil
}

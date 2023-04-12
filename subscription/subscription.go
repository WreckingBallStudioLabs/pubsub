package subscription

import (
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/common"
	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/name"
	"github.com/thalesfsp/configurer/util"
	"github.com/thalesfsp/status"
)

//////
// Const, vars, and types.
//////

// Func is the function to call when a message is received.
type Func func(msg *message.Message)

// Subscription is a subscription to a topic.
type Subscription struct {
	common.Common

	// Func is the function to call when a message is received. It handles the
	// message.
	Func Func `json:"-"`

	// Channel is the channel to receive messages.
	Channel chan *message.Message `json:"-"`
}

//////
// Factory.
//////

// New creates a new subscription. topic and queue should be in the form of the
// following example: "v1.meta.created" and "v1.meta.created.queue".
func New(topic, queue string, callback Func) (*Subscription, error) {
	t, err := name.New(topic)
	if err != nil {
		return nil, err
	}

	q, err := name.New(queue)
	if err != nil {
		return nil, err
	}

	s := &Subscription{
		Common: common.Common{
			CreatedAt: time.Now(),
			Status:    status.Created,
			Queue:     q.String(),
			Topic:     t.String(),
		},

		Func:    callback,
		Channel: make(chan *message.Message),
	}

	if err := util.Process(s); err != nil {
		return nil, err
	}

	return s, nil
}

// MustNew creates a new subscription, panicking if there's an error.
func MustNew(topic, queue string, callback Func) *Subscription {
	s, err := New(topic, queue, callback)
	if err != nil {
		panic(err)
	}

	return s
}

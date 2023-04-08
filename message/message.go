package message

import (
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/common"
	"github.com/WreckingBallStudioLabs/pubsub/name"
	"github.com/thalesfsp/configurer/util"
	"github.com/thalesfsp/status"
)

//////
// Const, vars, and types.
//////

// Message definition.
type Message struct {
	common.Common

	// Data to be published.
	Data any `json:"data" validate:"required" id:"uuid"`

	// Topic to publish to.
	Topic string `json:"topic"`
}

//////
// Factory.
//////

// New creates a new message.
func New(topic string, data any) (*Message, error) {
	t, err := name.New(topic)
	if err != nil {
		return nil, err
	}

	m := &Message{
		Common: common.Common{
			CreatedAt: time.Now(),
			Status:    status.Created,
		},

		Data:  data,
		Topic: t.String(),
	}

	if err := util.Process(m); err != nil {
		panic(err)
	}

	return m, nil
}

// MustNew creates a new message, panicking if there's an error.
func MustNew(topic string, data any) *Message {
	m, err := New(topic, data)
	if err != nil {
		panic(err)
	}

	return m
}

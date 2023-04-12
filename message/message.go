package message

import (
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/common"
	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
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
}

//////
// Methods.
//////

// Process the content of the message `b` into `v`.
func (m *Message) Process(b any, v any) error {
	jsonData, err := shared.Marshal(b)
	if err != nil {
		return err
	}

	if err := shared.Unmarshal(jsonData, &v); err != nil {
		return err
	}

	return nil
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
			Topic:     t.String(),
		},

		Data: data,
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

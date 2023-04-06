package nats

import (
	"context"

	"github.com/WreckingBallStudioLabs/pubsub/errorcatalog"
	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/WreckingBallStudioLabs/pubsub/pubsub"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	natsgo "github.com/nats-io/nats.go"
	"github.com/thalesfsp/customerror"
)

//////
// Const, vars, and types.
//////

// Name is the name of the pubsub.
const Name = "nats"

// Singleton.
var singleton pubsub.IPubSub

// Option is for the NATS configuration.
type Option = natsgo.Option

// NATS pubsub definition.
type NATS struct {
	*pubsub.PubSub

	// Options are the NATS configuration.
	Options []Option `json:"-" validate:"required"`

	// Client is the NATS client.
	Client *natsgo.Conn

	// URL is the NATS URL.
	URL string `json:"url" validate:"required"`
}

//////
// Helpers
//////

// MessageToPayload converts a message to a payload.
func MessageToPayload(message interface{}) ([]byte, error) {
	switch value := message.(type) {
	case []byte:
		return shared.Marshal(string(value))
	case string:
		return shared.Marshal(value)
	default:
		return shared.MarshalIndent(value, "", "  ")
	}
}

//////
// Implement the PubSubClient interface.
//////

// Publish sends a message to a topic.
func (c *NATS) Publish(topic string, message interface{}) error {
	payload, err := MessageToPayload(message)
	if err != nil {
		return err
	}

	if err := c.Client.Publish(topic, payload); err != nil {
		return errorcatalog.
			Get().
			MustGet(
				errorcatalog.PubSubErrNatsPublish,
				customerror.WithError(err),
				customerror.WithField("topic", topic),
			)
	}

	return nil
}

// PublishAsync sends a message to a topic. In case of error it will just log
// it.
func (c *NATS) PublishAsync(topic string, message any) {
	go func() {
		if err := c.Publish(topic, message); err != nil {
			c.GetLogger().Error(err)
		}
	}()
}

// Subscribe to a topic.
func (c *NATS) Subscribe(topic string, queue string, cb func([]byte)) (subscription.Subscription, error) {
	ch := make(chan []byte)
	sub := subscription.Subscription{
		Topic:    topic,
		Queue:    queue,
		Callback: cb,
		Channel:  ch, // Use a receive-only channel for subscriptions
	}

	_, err := c.Client.QueueSubscribe(topic, queue, func(m *natsgo.Msg) {
		ch <- m.Data
	})
	if err != nil {
		close(ch)

		sub.Channel = nil

		return sub, errorcatalog.
			Get().
			MustGet(
				errorcatalog.PubSubErrNatsSubscribe,
				customerror.WithError(err),
				customerror.WithField("topic", topic),
			)
	}

	return sub, nil
}

// Unsubscribe from a topic.
func (c *NATS) Unsubscribe(topic string) error {
	return nil
}

// Close the connection to the Pub Sub broker.
func (c *NATS) Close() error {
	c.Client.Close()

	return nil
}

// GetClient returns the storage client. Use that to interact with the
// underlying storage client.
func (c *NATS) GetClient() any {
	return c.Client
}

//////
// Factory.
//////

// New creates a new NATS pubsub.
func New(ctx context.Context, url string, options ...Option) (pubsub.IPubSub, error) {
	var _ pubsub.IPubSub = (*NATS)(nil)

	p, err := pubsub.New(Name)
	if err != nil {
		return nil, err
	}

	natsConn, err := natsgo.Connect(url, options...)
	if err != nil {
		return nil, err
	}

	client := &NATS{
		PubSub: p,

		Client:  natsConn,
		Options: options,
		URL:     url,
	}

	singleton = client

	return client, nil
}

//////
// Exported functionalities.
//////

// Get returns a setup NATS, or set it up.
func Get() pubsub.IPubSub {
	if singleton == nil {
		panic(errorcatalog.Get().MustGet(errorcatalog.PubSubErrNatsNilMessage))
	}

	return singleton
}

// Set sets the singleton. Useful for testing.
func Set(ps pubsub.IPubSub) {
	singleton = ps
}

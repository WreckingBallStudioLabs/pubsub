package nats

import (
	"context"

	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/WreckingBallStudioLabs/pubsub/pubsub"
	natsgo "github.com/nats-io/nats.go"
)

//////
// Const, vars, and types.
//////

// Singleton.
var singleton pubsub.IPubSub

// Config is for the NATS configuration.
type Config = natsgo.Options

// NATS pubsub definition.
type NATS struct {
	*pubsub.PubSub

	// Config is the NATS configuration.
	Config *Config `json:"-" validate:"required"`

	// Client is the NATS client.
	Client *natsgo.Conn
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

	return c.Client.Publish(topic, payload)
}

// Subscribe subscribes to a topic and returns a channel for receiving messages.
func (c *NATS) Subscribe(topic string, queue string, cb func([]byte)) pubsub.Subscription {
	ch := make(chan []byte)
	sub := pubsub.Subscription{
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
	}

	return sub
}

// Unsubscribe unsubscribes from a topic.
func (c *NATS) Unsubscribe(topic string) error {
	return nil
}

// Close closes the connection to the Pub Sub broker.
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
func New(ctx context.Context, cfg Config) (*NATS, error) {
	var _ pubsub.IPubSub = (*NATS)(nil)

	natsConn, err := natsgo.Connect(
		cfg.Url,
		natsgo.MaxReconnects(cfg.MaxReconnect),
		natsgo.ReconnectWait(cfg.ReconnectWait),
		natsgo.Timeout(cfg.Timeout),
		natsgo.PingInterval(cfg.PingInterval),
		natsgo.Name(cfg.Name),
	)
	if err != nil {
		return nil, err
	}

	client := &NATS{
		Client: natsConn,
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
		panic("NATS client not initialized")
	}

	return singleton
}

// Set sets the singleton. Useful for testing.
func Set(ps pubsub.IPubSub) {
	singleton = ps
}

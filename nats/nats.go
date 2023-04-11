package nats

import (
	"context"

	"github.com/WreckingBallStudioLabs/pubsub/errorcatalog"
	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/pubsub"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	natsgo "github.com/nats-io/nats.go"
	"github.com/thalesfsp/concurrentloop"
	"github.com/thalesfsp/customerror"
	"github.com/thalesfsp/validation"
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
// Implement the PubSubClient interface.
//////

// Publish sends a message to a topic.
func (c *NATS) Publish(
	ctx context.Context,
	messages []*message.Message,
	opts ...pubsub.Func,
) ([]*message.Message, concurrentloop.Errors) {
	r, err := concurrentloop.Map(
		ctx, messages,
		func(ctx context.Context, message *message.Message) (*message.Message, error) {
			if err := validation.Validate(message); err != nil {
				return message, err
			}

			// Process options.
			o, err := pubsub.NewOptions()
			if err != nil {
				return message, err
			}

			for _, opt := range opts {
				if err := opt(o); err != nil {
					return message, err
				}
			}

			if o.Sync {
				return message, errorcatalog.
					Get().
					MustGet(errorcatalog.PubSubErrPubSubNotImpl).
					NewFailedToError(
						customerror.WithField("topic", message.Topic),
						customerror.WithField("id", message.ID),
					)
			}

			payload, err := shared.MarshalIndent(message, "", "  ")
			if err != nil {
				return message, err
			}

			if err := c.Client.Publish(message.Topic, payload); err != nil {
				return message, errorcatalog.
					Get().
					MustGet(
						errorcatalog.PubSubErrNATSPublish,
						customerror.WithError(err),
						customerror.WithField("topic", message.Topic),
						customerror.WithField("id", message.ID),
					).NewFailedToError()
			}

			return message, nil
		})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// MustPublish sends a message to a topic. In case of error it will panic.
func (c *NATS) MustPublish(ctx context.Context, msgs ...*message.Message) []*message.Message {
	messages, err := c.Publish(ctx, msgs)
	if err != nil {
		panic(err)
	}

	return messages
}

// MustPublishAsync sends a message to a topic asynchronously. In case of error
// it will panic.
func (c *NATS) MustPublishAsync(ctx context.Context, messages ...*message.Message) {
	go c.MustPublish(ctx, messages...)
}

// Subscribe to a topic.
func (c *NATS) Subscribe(ctx context.Context, subscriptions ...*subscription.Subscription) ([]*subscription.Subscription, concurrentloop.Errors) {
	r, err := concurrentloop.Map(ctx, subscriptions, func(ctx context.Context, subscription *subscription.Subscription) (*subscription.Subscription, error) {
		if err := validation.Validate(subscription); err != nil {
			return subscription, err
		}

		_, err := c.Client.QueueSubscribe(subscription.Topic, subscription.Queue, func(m *natsgo.Msg) {
			var msg message.Message

			if err := shared.Unmarshal(m.Data, &msg); err != nil {
				// TODO: Handle this error with APM, metrics, etc.
				panic(err)
			}

			// Runs the subscription handler function.
			subscription.Func(&msg)

			// Also sends the data to the channel.
			subscription.Channel <- &msg
		})
		if err != nil {
			close(subscription.Channel)

			subscription.Channel = nil

			return subscription, errorcatalog.
				Get().
				MustGet(
					errorcatalog.PubSubErrNATSSubscribe,
					customerror.WithError(err),
					customerror.WithField("topic", subscription.Topic),
					customerror.WithField("id", subscription.ID),
				).NewFailedToError()
		}

		return subscription, nil
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// MustSubscribe to a topic. In case of error it will panic.
func (c *NATS) MustSubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) []*subscription.Subscription {
	subscriptions, err := c.Subscribe(ctx, subscriptions...)
	if err != nil {
		panic(err)
	}

	return subscriptions
}

// MustSubscribeAsyn to a topic asynchronously. In case of error it will panic.
func (c *NATS) MustSubscribeAsyn(ctx context.Context, subscriptions ...*subscription.Subscription) {
	go c.MustSubscribe(ctx, subscriptions...)
}

// Unsubscribe from a topic.
func (c *NATS) Unsubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) error {
	c.GetLogger().Warnln(
		errorcatalog.Get().MustGet(errorcatalog.PubSubErrPubSubNotImpl).NewFailedToError(),
	)

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
		panic(errorcatalog.Get().MustGet(errorcatalog.PubSubErrNATANilMessage).NewFailedToError())
	}

	return singleton
}

// Set sets the singleton. Useful for testing.
func Set(ps pubsub.IPubSub) {
	singleton = ps
}

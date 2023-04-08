package pubsub

import (
	"context"

	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	"github.com/thalesfsp/concurrentloop"
	"github.com/thalesfsp/sypl"
)

//////
// Vars, consts, and types.
//////

// IPubSub defines a PubSub does.
type IPubSub interface {
	// Publish sends a message to a topic.
	Publish(ctx context.Context, messages ...*message.Message) ([]*message.Message, concurrentloop.Errors)

	// MustPublish sends a message to a topic. In case of error it will panic.
	MustPublish(ctx context.Context, msgs ...*message.Message) []*message.Message

	// MustPublishAsync sends a message to a topic asynchronously. In case of
	// error it will panic.
	MustPublishAsync(ctx context.Context, messages ...*message.Message)

	// Subscribe to a topic.
	Subscribe(ctx context.Context, v any, subscriptions ...*subscription.Subscription) ([]*subscription.Subscription, concurrentloop.Errors)

	// MustSubscribe to a topic. In case of error it will panic.
	MustSubscribe(ctx context.Context, v any, subscriptions ...*subscription.Subscription) []*subscription.Subscription

	// MustSubscribeAsyn to a topic asynchronously. In case of error it will panic.
	MustSubscribeAsyn(ctx context.Context, v any, subscriptions ...*subscription.Subscription)

	// Unsubscribe from a topic.
	Unsubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) error

	// Close the connection to the Pub Sub broker.
	Close() error

	// GetClient returns the storage client. Use that to interact with the
	// underlying storage client.
	GetClient() any

	// GetLogger returns the logger.
	GetLogger() sypl.ISypl

	// GetName returns the pubsub name.
	GetName() string
}

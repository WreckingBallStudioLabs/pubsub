package pubsub

import (
	"context"
	"expvar"

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
	Publish(ctx context.Context, messages []*message.Message, opts ...Func) ([]*message.Message, concurrentloop.Errors)

	// MustPublish sends a message to a topic. In case of error it will panic.
	MustPublish(ctx context.Context, msgs ...*message.Message) []*message.Message

	// MustPublishAsync sends a message to a topic asynchronously. In case of
	// error it will panic.
	MustPublishAsync(ctx context.Context, messages ...*message.Message)

	// Subscribe to a topic.
	Subscribe(ctx context.Context, subscriptions []*subscription.Subscription, opts ...Func) ([]*subscription.Subscription, concurrentloop.Errors)

	// MustSubscribe to a topic. In case of error it will panic.
	MustSubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) []*subscription.Subscription

	// MustSubscribeAsyn to a topic asynchronously. In case of error it will panic.
	MustSubscribeAsyn(ctx context.Context, subscriptions ...*subscription.Subscription)

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

	// GetType returns its type.
	GetType() string

	// GetCounterPingFailed returns the metric.
	GetCounterPingFailed() *expvar.Int

	// GetPublishedCounter returns the metric.
	GetPublishedCounter() *expvar.Int

	// GetPublishedFailedCounter returns the metric.
	GetPublishedFailedCounter() *expvar.Int

	// GetSubscribedCounter returns the metric.
	GetSubscribedCounter() *expvar.Int

	// GetSubscribedFailedCounter returns the metric.
	GetSubscribedFailedCounter() *expvar.Int
}

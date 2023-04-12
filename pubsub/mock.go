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
// Creates the a struct which satisfies the storage.IStorage interface.
//////

// Mock is a struct which satisfies the pubsub.IPubSub interface.
//
//nolint:dupl
type Mock struct {
	// Publish sends a message to a topic.
	MockPublish func(ctx context.Context, messages []*message.Message, opts ...Func) ([]*message.Message, concurrentloop.Errors)

	// MustPublish sends a message to a topic. In case of error it will panic.
	MockMustPublish func(ctx context.Context, msgs ...*message.Message) []*message.Message

	// MustPublishAsync sends a message to a topic asynchronously. In case of error it will panic.
	MockMustPublishAsync func(ctx context.Context, messages ...*message.Message)

	// Subscribe to a topic.
	MockSubscribe func(ctx context.Context, subscriptions []*subscription.Subscription, opts ...Func) ([]*subscription.Subscription, concurrentloop.Errors)

	// MustSubscribe to a topic. In case of error it will panic.
	MockMustSubscribe func(ctx context.Context, subscriptions ...*subscription.Subscription) []*subscription.Subscription

	// MustSubscribeAsync to a topic asynchronously. In case of error it will panic.
	MockMustSubscribeAsyn func(ctx context.Context, subscriptions ...*subscription.Subscription)

	// Unsubscribe from a topic.
	MockUnsubscribe func(ctx context.Context, subscriptions ...*subscription.Subscription) error

	// Close the connection to the Pub Sub broker.
	MockClose func() error

	// GetClient returns the storage client. Use that to interact with the underlying storage client.
	MockGetClient func() any

	// GetLogger returns the logger.
	MockGetLogger func() sypl.ISypl

	// GetName returns the pubsub name.
	MockGetName func() string

	// GetType returns its type.
	MockGetType func() string

	// GetCounterPingFailed returns the metric.
	MockGetCounterPingFailed func() *expvar.Int

	// GetPublishedCounter returns the metric.
	MockGetPublishedCounter func() *expvar.Int

	// GetPublishedFailedCounter returns the metric.
	MockGetPublishedFailedCounter func() *expvar.Int

	// GetSubscribedCounter returns the metric.
	MockGetSubscribedCounter func() *expvar.Int

	// GetSubscribedFailedCounter returns the metric.
	MockGetSubscribedFailedCounter func() *expvar.Int
}

//////
// When the methods are called, it will call the corresponding method in the
// Mock struct returning the desired value. This implements the IStorage
// interface.
//////

// Publish sends a message to a topic.
func (m *Mock) Publish(ctx context.Context, messages []*message.Message, opts ...Func) ([]*message.Message, concurrentloop.Errors) {
	return m.MockPublish(ctx, messages, opts...)
}

// MustPublish sends a message to a topic. In case of error it will panic.
func (m *Mock) MustPublish(ctx context.Context, msgs ...*message.Message) []*message.Message {
	return m.MockMustPublish(ctx, msgs...)
}

// MustPublishAsync sends a message to a topic asynchronously. In case of error it will panic.
func (m *Mock) MustPublishAsync(ctx context.Context, messages ...*message.Message) {
	m.MockMustPublishAsync(ctx, messages...)
}

// Subscribe to a topic.
func (m *Mock) Subscribe(ctx context.Context, subscriptions []*subscription.Subscription, opts ...Func) ([]*subscription.Subscription, concurrentloop.Errors) {
	return m.MockSubscribe(ctx, subscriptions, opts...)
}

// MustSubscribe to a topic. In case of error it will panic.
func (m *Mock) MustSubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) []*subscription.Subscription {
	return m.MockMustSubscribe(ctx, subscriptions...)
}

// MustSubscribeAsync to a topic asynchronously. In case of error it will panic.
func (m *Mock) MustSubscribeAsync(ctx context.Context, subscriptions ...*subscription.Subscription) {
	m.MockMustSubscribeAsyn(ctx, subscriptions...)
}

// Unsubscribe from a topic.
func (m *Mock) Unsubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) error {
	return m.MockUnsubscribe(ctx, subscriptions...)
}

// Close the connection to the Pub Sub broker.
func (m *Mock) Close() error {
	return m.MockClose()
}

// GetClient returns the storage client. Use that to interact with the underlying storage client.
func (m *Mock) GetClient() any {
	return m.MockGetClient()
}

// GetLogger returns the logger.
func (m *Mock) GetLogger() sypl.ISypl {
	return m.MockGetLogger()
}

// GetName returns the pubsub name.
func (m *Mock) GetName() string {
	return m.MockGetName()
}

// GetType returns its type.
func (m *Mock) GetType() string {
	return m.MockGetType()
}

// GetCounterPingFailed returns the metric.
func (m *Mock) GetCounterPingFailed() *expvar.Int {
	return m.MockGetCounterPingFailed()
}

// GetPublishedCounter returns the metric.
func (m *Mock) GetPublishedCounter() *expvar.Int {
	return m.MockGetPublishedCounter()
}

// GetPublishedFailedCounter returns the metric.
func (m *Mock) GetPublishedFailedCounter() *expvar.Int {
	return m.MockGetPublishedFailedCounter()
}

// GetSubscribedCounter returns the metric.
func (m *Mock) GetSubscribedCounter() *expvar.Int {
	return m.MockGetSubscribedCounter()
}

// GetSubscribedFailedCounter returns the metric.
func (m *Mock) GetSubscribedFailedCounter() *expvar.Int {
	return m.MockGetSubscribedFailedCounter()
}

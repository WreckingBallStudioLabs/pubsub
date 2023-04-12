package pubsub

import (
	"context"

	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	"github.com/thalesfsp/concurrentloop"
)

//////
// Vars, consts, and types.
//////

// Map is a map of PubSubs
type Map map[string]IPubSub

// PublishMany will make all PubSubs to concurrently publish many messages.
func (m Map) PublishMany(
	ctx context.Context,
	messages []*message.Message,
	opts ...Func,
) ([]*message.Message, error) {
	var (
		msgs []*message.Message
		errs concurrentloop.Errors
	)

	for _, pubsub := range m {
		m, e := pubsub.Publish(ctx, messages, opts...)
		msgs = append(msgs, m...)
		errs = append(errs, e)
	}

	return msgs, errs
}

// MustPublishManyAsync will make all PubSubs to concurrently publish many messages
// asynchronously.
func (m Map) MustPublishManyAsync(ctx context.Context, messages ...*message.Message) {
	go func() {
		for _, pubsub := range m {
			pubsub.MustPublishAsync(ctx, messages...)
		}
	}()
}

// SubscribeMany will make all PubSubs to concurrently subscribe to many
// subscriptions.
func (m Map) SubscribeMany(
	ctx context.Context,
	subscriptions []*subscription.Subscription,
	opts ...Func,
) ([]*subscription.Subscription, concurrentloop.Errors) {
	var (
		msgs []*subscription.Subscription
		errs concurrentloop.Errors
	)

	for _, pubsub := range m {
		m, e := pubsub.Subscribe(ctx, subscriptions, opts...)
		msgs = append(msgs, m...)
		errs = append(errs, e)
	}

	return msgs, errs
}

// MustSubscribeManyAsync will make all PubSubs to concurrently subscribe to many
// subscriptions asynchronously.
func (m Map) MustSubscribeManyAsync(ctx context.Context, subscriptions ...*subscription.Subscription) {
	go func() {
		for _, pubsub := range m {
			pubsub.MustSubscribeAsyn(ctx, subscriptions...)
		}
	}()
}

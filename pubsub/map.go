package pubsub

import (
	"context"

	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
)

//////
// Vars, consts, and types.
//////

// Map is a map of PubSubs
type Map map[string]IPubSub

// PublishMany will make all PubSubs to concurrently publish many messages.
func (m Map) PublishMany(ctx context.Context, messages ...*message.Message) {
	for _, pubsub := range m {
		pubsub.Publish(ctx, messages...)
	}
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
func (m Map) SubscribeMany(ctx context.Context, v any, subscriptions ...*subscription.Subscription) {
	for _, pubsub := range m {
		pubsub.Subscribe(ctx, v, subscriptions...)
	}
}

// MustSubscribeManyAsync will make all PubSubs to concurrently subscribe to many
// subscriptions asynchronously.
func (m Map) MustSubscribeManyAsync(ctx context.Context, v any, subscriptions ...*subscription.Subscription) {
	go func() {
		for _, pubsub := range m {
			pubsub.MustSubscribeAsyn(ctx, v, subscriptions...)
		}
	}()
}

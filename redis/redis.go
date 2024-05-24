package redis

import (
	"context"
	"sync"
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/errorcatalog"
	"github.com/WreckingBallStudioLabs/pubsub/internal/customapm"
	"github.com/WreckingBallStudioLabs/pubsub/internal/logging"
	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/pubsub"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	"github.com/eapache/go-resiliency/retrier"
	redisgo "github.com/redis/go-redis/v9"
	"github.com/thalesfsp/concurrentloop"
	"github.com/thalesfsp/customerror"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/fields"
	"github.com/thalesfsp/sypl/level"
	"github.com/thalesfsp/validation"
)

//////
// Const, vars, and types.
//////

// Name is the name of the pubsub.
const Name = "redisgo"

// Singleton.
var singleton pubsub.IPubSub

// Option is for the Redis configuration.
type Option = redisgo.Options

// Redis pubsub definition.
type Redis struct {
	*pubsub.PubSub

	// Options are the Redis configuration.
	Options []Option `json:"-" validate:"required"`

	// Client is the Redis client.
	Client *redisgo.Client

	// URL is the Redis URL.
	URL string `json:"url" validate:"required"`

	// Mutex to synchronize access to subscription channels.
	mu sync.RWMutex
}

//////
// Implement the PubSubClient interface.
//////

// Publish sends a message to a topic.
func (r *Redis) Publish(
	ctx context.Context,
	messages []*message.Message,
	opts ...pubsub.Func,
) ([]*message.Message, concurrentloop.Errors) {
	// APM Tracing.
	ctx, span := customapm.Trace(
		ctx,
		r.GetType(),
		Name,
		status.Published.String(),
	)
	defer span.End()

	// Publish messages concurrently.
	result, errs := concurrentloop.Map(
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
				return message, customapm.TraceError(
					ctx,
					errorcatalog.
						Get().
						MustGet(errorcatalog.PubSubErrPubSubNotImpl).
						NewFailedToError(
							customerror.WithField("topic", message.Topic),
							customerror.WithField("id", message.ID),
						),
					r.GetLogger(),
					r.GetPublishedFailedCounter(),
				)
			}

			// Serialize message payload.
			payload, err := shared.Marshal(message)
			if err != nil {
				return message, err
			}

			// Publish the message to Redis.
			if err := r.Client.Publish(ctx, message.Topic, payload).Err(); err != nil {
				return message, errorcatalog.
					Get().
					MustGet(
						errorcatalog.PubSubErrPublish,
						customerror.WithError(err),
						customerror.WithField("topic", message.Topic),
						customerror.WithField("id", message.ID),
					).NewFailedToError()
			}

			return message, nil
		})
	if len(errs) > 0 {
		_ = customapm.TraceError(ctx, errs, r.GetLogger(), r.GetPublishedFailedCounter())

		return nil, errs
	}

	// Logging
	r.GetLogger().PrintlnWithOptions(
		level.Debug,
		status.Published.String(),
		sypl.WithFields(logging.ToAPM(ctx, make(fields.Fields))),
	)

	// Metrics.
	r.GetPublishedCounter().Add(1)

	return result, nil
}

// MustPublish sends a message to a topic. In case of error it will panic.
func (r *Redis) MustPublish(ctx context.Context, msgs ...*message.Message) []*message.Message {
	messages, err := r.Publish(ctx, msgs)
	if err != nil {
		panic(err)
	}

	return messages
}

// MustPublishAsync sends a message to a topic asynchronously. In case of error it will panic.
func (r *Redis) MustPublishAsync(ctx context.Context, messages ...*message.Message) {
	go r.MustPublish(ctx, messages...)
}

// Subscribe to a topic.
func (r *Redis) Subscribe(
	ctx context.Context,
	subscriptions []*subscription.Subscription,
	opts ...pubsub.Func,
) ([]*subscription.Subscription, concurrentloop.Errors) {
	// APM Tracing.
	ctx, span := customapm.Trace(
		ctx,
		r.GetType(),
		Name,
		status.Subscribed.String(),
	)
	defer span.End()

	//////
	// Subscribe concurrently.
	//////

	result, err := concurrentloop.Map(
		ctx,
		subscriptions,
		func(ctx context.Context, subscription *subscription.Subscription) (*subscription.Subscription, error) {
			if err := validation.Validate(subscription); err != nil {
				return subscription, err
			}

			//////
			// Process options.
			//////

			o, err := pubsub.NewOptions()
			if err != nil {
				return subscription, err
			}

			for _, opt := range opts {
				if err := opt(o); err != nil {
					return subscription, err
				}
			}

			if o.Sync {
				return subscription, errorcatalog.
					Get().
					MustGet(errorcatalog.PubSubErrPubSubNotImpl).
					NewFailedToError(
						customerror.WithField("topic", subscription.Topic),
						customerror.WithField("id", subscription.ID),
					)
			}

			pubsub := r.Client.Subscribe(ctx, subscription.Topic)

			go func() {
				defer close(subscription.Channel)

				ch := pubsub.Channel()

				for msg := range ch {
					var m message.Message

					if err := shared.Unmarshal([]byte(msg.Payload), &m); err != nil {
						close(subscription.Channel)

						subscription.Channel = nil

						panic(customapm.TraceError(ctx, err, r.GetLogger(), r.GetSubscribedFailedCounter()))
					}

					// Runs the subscription handler function.
					subscription.Func(&m)

					// Also sends the data to the channel.
					subscription.Channel <- &m
				}
			}()

			return subscription, nil
		})
	if err != nil {
		_ = customapm.TraceError(ctx, err, r.GetLogger(), r.GetSubscribedFailedCounter())

		return nil, err
	}

	// Logging
	r.GetLogger().PrintlnWithOptions(
		level.Debug,
		status.Subscribed.String(),
		sypl.WithFields(logging.ToAPM(ctx, make(fields.Fields))),
	)

	// Metrics.
	r.GetSubscribedCounter().Add(1)

	return result, nil
}

// MustSubscribe to a topic. In case of error it will panic.
func (r *Redis) MustSubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) []*subscription.Subscription {
	subscriptions, err := r.Subscribe(ctx, subscriptions)
	if err != nil {
		panic(err)
	}

	return subscriptions
}

// MustSubscribeAsync to a topic asynchronously. In case of error it will panic.
func (r *Redis) MustSubscribeAsync(ctx context.Context, subscriptions ...*subscription.Subscription) {
	go r.MustSubscribe(ctx, subscriptions...)
}

// Unsubscribe from a topic.
func (r *Redis) Unsubscribe(ctx context.Context, subscriptions ...*subscription.Subscription) error {
	return customapm.TraceError(
		ctx,
		errorcatalog.
			Get().
			MustGet(errorcatalog.PubSubErrPubSubNotImpl).
			NewFailedToError(),
		r.GetLogger(),
		r.GetPublishedFailedCounter(),
	)
}

// Close the connection to the Pub Sub broker.
func (r *Redis) Close() error {
	return r.Client.Close()
}

// GetClient returns the storage client. Use that to interact with the underlying storage client.
func (r *Redis) GetClient() any {
	return r.Client
}

//////
// Factory.
//////

// New creates a new Redis pubsub.
func New(ctx context.Context, url string, options ...Option) (pubsub.IPubSub, error) {
	var _ pubsub.IPubSub = (*Redis)(nil)

	p, err := pubsub.New(ctx, Name)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = make([]Option, 1)
	}

	opt := options[0]

	opt.Addr = url

	client := redisgo.NewClient(&opt)

	r := retrier.New(retrier.ExponentialBackoff(3, 10*time.Second), nil)

	if err := r.Run(func() error {
		if err := client.Ping(ctx).Err(); err != nil {
			return customerror.NewFailedToError("ping", customerror.WithError(err))
		}

		return nil
	}); err != nil {
		return nil, customapm.TraceError(ctx, err, p.GetLogger(), p.GetCounterPingFailed())
	}

	ps := &Redis{
		PubSub:  p,
		Options: options,
		Client:  client,
		URL:     url,
		mu:      sync.RWMutex{},
	}

	singleton = ps

	return ps, nil
}

//////
// Exported functionalities.
//////

// Get returns a setup Redis, or set it up.
func Get() pubsub.IPubSub {
	if singleton == nil {
		panic(errorcatalog.Get().MustGet(errorcatalog.PubSubErrNilClient).NewFailedToError())
	}

	return singleton
}

// Set sets the singleton. Useful for testing.
func Set(ps pubsub.IPubSub) {
	singleton = ps
}

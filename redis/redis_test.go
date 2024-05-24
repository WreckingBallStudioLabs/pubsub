package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/errorcatalog"
	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	if !shared.IsEnvironment(shared.Integration) {
		t.Skip("Skipping test. Not in e2e " + shared.Integration + "environment.")
	}

	t.Setenv("PUBSUB_METRICS_PREFIX", "redis_test")

	host := os.Getenv("REDIS_HOST")

	if host == "" {
		t.Fatal("REDIS_HOST is not set")
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Should work - E2E",
			args: args{
				ctx: context.Background(),
				id:  shared.DocumentID,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tear up.
			ctx, cancel := context.WithTimeout(tt.args.ctx, shared.DefaultTimeout)
			defer cancel()

			client, err := New(ctx, host)
			assert.NoError(t, err)
			assert.NotNil(t, client)

			if client != nil && client.GetClient() == nil {
				t.Fatal("client.Client is nil")
			}

			// Create a subscription with a callback function.
			sub := subscription.MustNew("v1.meta.created", "v1.meta.created.queue", func(msg *message.Message) {
				var v shared.TestDataS

				if err := msg.Process(msg.Data, &v); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, shared.TestData, &v)
			})

			// And here is the channel way.
			//
			// NOTE: The channel is buffered, so it's important to read from it
			// in a goroutine.
			// NOTE: Optionally, listen to the `ctx.Done` channel to stop the
			// goroutine.
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case msg := <-sub.Channel:
						// NOTE: Need to handle if msg is nil because the channel
						// is buffered.
						if msg != nil {
							var v shared.TestDataS

							if err := msg.Process(msg.Data, &v); err != nil {
								panic(err)
							}

							assert.Equal(t, shared.TestData, &v)
						}
					}
				}
			}()

			// Subscribe to the channel.
			assert.NotPanics(t, func() {
				client.MustSubscribe(ctx, sub)
			})

			// Publish to the channel.
			if sub != nil && sub.Topic != "" {
				assert.NotPanics(t, func() {
					client.MustPublish(ctx, message.MustNew(sub.Topic, shared.TestData))
				})
			}

			// Wait for the message to be processed.
			time.Sleep(5 * time.Second)

			// Check if the metrics are working.
			assert.Equal(t, int64(1), client.GetPublishedCounter().Value())
			assert.Equal(t, int64(0), client.GetPublishedFailedCounter().Value())
			assert.Equal(t, int64(1), client.GetSubscribedCounter().Value())
			assert.Equal(t, int64(0), client.GetSubscribedFailedCounter().Value())

			// Unsubscribe from the topic.
			ErrPubSubPubSubNotImpl := errorcatalog.Get().MustGet(errorcatalog.PubSubErrPubSubNotImpl).NewFailedToError()
			assert.ErrorAs(t, client.Unsubscribe(ctx, sub), &ErrPubSubPubSubNotImpl)
		})
	}
}

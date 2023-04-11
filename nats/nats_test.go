package nats

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// if !shared.IsEnvironment(shared.Integration) {
	// 	t.Skip("Skipping test. Not in e2e " + shared.Integration + "environment.")
	// }

	t.Setenv("PUBSUB_METRICS_PREFIX", "test")

	host := os.Getenv("NATS_HOST")

	if host == "" {
		t.Fatal("NATS_HOST is not set")
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
			name: "Shoud work - E2E",
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
			//////
			// Tear up.
			//////

			ctx, cancel := context.WithTimeout(tt.args.ctx, shared.DefaultTimeout)
			defer cancel()

			client, err := New(ctx, host)
			assert.NoError(t, err)
			assert.NotNil(t, client)

			if client != nil && client.GetClient() == nil {
				t.Fatal("client.Client is nil")
			}

			//////
			// Should be able to subscribe to a channel.
			//////

			// 1. Create a subscription.
			//
			// Data can be handled in a channel or a callback. Here is the
			// callback way.
			//
			// NOTE: the topic and queue pattern, it's enforced by validation.
			// NOTE: The only moment something is typed, is the topic and queue.
			sub := subscription.MustNew("v1.meta.created", "v1.meta.created.queue", func(msg *message.Message) {
				var v shared.TestDataS

				// `Message` provides `Process` to easily unmarshal the data.
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
						var v shared.TestDataS

						if err := msg.Process(msg.Data, &v); err != nil {
							panic(err)
						}

						assert.Equal(t, shared.TestData, &v)
					}
				}
			}()

			// 2. Subscribe to the channel.
			subs, err := client.Subscribe(ctx, sub)
			assert.Nil(t, err)
			assert.NotNil(t, subs)

			//////
			// Should be able to publish to a channel.
			//////

			if subs != nil && len(subs) > 0 {
				assert.NotPanics(t, func() {
					// 3. Create a message (`MustNew`) then...
					//
					// 4. Publish to the channel.
					//
					// Smartly reuse `subs` topic, less typing, less error prone.
					client.MustPublish(ctx, message.MustNew(subs[0].Topic, shared.TestData))
				})
			}

			// Need to wait for the message to be processed.
			time.Sleep(5 * time.Second)
		})
	}
}

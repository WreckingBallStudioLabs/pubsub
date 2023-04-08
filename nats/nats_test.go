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

			//////
			// Should be able to subscribe to a channel.
			//////

			sub := subscription.MustNew("v1.meta.created", "v1.meta.created.queue", func(msg *message.Message) {
				t.Logf("HERE: %+v", msg.Data)
			})

			// Print messages from the channel.
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case msg := <-sub.Channel:
						t.Logf("HERE2: %+v", msg.Data)
					}
				}
			}()

			var v *shared.TestDataS
			subs, err := client.Subscribe(ctx, v, sub)
			t.Logf("HERE3: %+v", v)
			assert.Nil(t, err)

			//////
			// Should be able to publish to a channel.
			//////

			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, "test"))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, 1))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, 1.1))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, true))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, false))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, []int{1, 2, 3}))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, []string{"a", "b", "c"}))
			// _, err = client.Publish(ctx, message.MustNew(subs[0].Topic, map[int]string{1: "a", 2: "b", 3: "c"}))
			if subs != nil && len(subs) > 0 {
				_, err = client.Publish(ctx, message.MustNew(subs[0].Topic, shared.TestData))

				assert.Nil(t, err)
			}

			time.Sleep(5 * time.Second)
		})
	}
}

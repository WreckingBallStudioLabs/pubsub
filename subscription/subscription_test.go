package subscription

import (
	"testing"
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/common"
	"github.com/WreckingBallStudioLabs/pubsub/message"
	"github.com/stretchr/testify/assert"
	"github.com/thalesfsp/status"
)

func TestNew(t *testing.T) {
	type args struct {
		topic    string
		queue    string
		callback Func
	}
	tests := []struct {
		name    string
		args    args
		want    *Subscription
		wantErr bool
	}{
		{
			name: "Should work",
			args: args{
				topic:    "v1.meta.created",
				queue:    "v1.meta.created.queue",
				callback: func(msg *message.Message) {},
			},
			want: &Subscription{
				Common: common.Common{
					CreatedAt: time.Now(),
					Status:    status.Created,
					Queue:     "v1.meta.created.queue",
					Topic:     "v1.meta.created",
				},
				Func:    func(msg *message.Message) {},
				Channel: make(chan *message.Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustNew(tt.args.topic, tt.args.queue, tt.args.callback)

			// Automatically generated fields.
			assert.NotEmpty(t, got.ID)
			assert.NotZero(t, got.CreatedAt)
			assert.Equal(t, status.Created, got.Status)

			// Manually set fields.
			assert.Equal(t, tt.args.topic, got.Topic)
			assert.Equal(t, tt.args.queue, got.Queue)
			assert.NotNil(t, tt.args.callback, got.Func)
			assert.NotNil(t, got.Channel)
		})
	}
}

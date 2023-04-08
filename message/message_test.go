package message

import (
	"testing"
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/common"
	"github.com/stretchr/testify/assert"
	"github.com/thalesfsp/status"
)

func TestNew(t *testing.T) {
	type args struct {
		topic string
		data  any
	}
	tests := []struct {
		name    string
		args    args
		want    *Message
		wantErr bool
	}{
		{
			name: "Should work",
			args: args{
				topic: "test",
				data:  "test",
			},
			want: &Message{
				Common: common.Common{
					CreatedAt: time.Now(),
					Status:    status.Created,
				},
				Topic: "test",
				Data:  "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustNew(tt.args.topic, tt.args.data)

			// Automatically generated fields.
			assert.NotEmpty(t, got.ID)
			assert.NotZero(t, got.CreatedAt)
			assert.Equal(t, status.Created, got.Status)

			// Manually set fields.
			assert.Equal(t, tt.args.topic, got.Topic)
			assert.Equal(t, tt.args.data, got.Data)
		})
	}
}

package message

import (
	"testing"

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
				topic: "v1.meta.created",
				data:  "v1.meta.created.queue",
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

package nats

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	if !shared.IsEnvironment(shared.Integration) {
		t.Skip("Skipping test. Not in e2e " + shared.Integration + "environment.")
	}

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

			_, err = client.Subscribe("test", "test-queue", func(msg []byte) {})
			assert.NoError(t, err)

			//////
			// Should be able to publish to a channel.
			//////

			assert.NoError(t, client.Publish("test", []byte("test")))

			//////
			// Should be able to publish using an interface.
			//////

			assert.NoError(t, client.Publish("test", struct {
				Test string
			}{
				Test: "test",
			}))

			//////
			// Should be able to publish using a map.
			//////

			assert.NoError(t, client.Publish("test", map[string]interface{}{
				"test": "test",
			}))

			//////
			// Should be able to publish using a byte array.
			//////

			assert.NoError(t, client.Publish("test", []byte("test")))
		})
	}
}

func TestInterfaceToBytes_String(t *testing.T) {
	input := "Hello, world!"

	expected := []byte(`"Hello, world!"`)

	result, err := MessageToPayload(input)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if !bytes.Equal(expected, result) {
		t.Errorf("Unexpected result: expected %v, got %v", expected, result)
	}
}

func TestInterfaceToBytes_ByteSlice(t *testing.T) {
	input := []byte{1, 2, 3, 4}

	expected := []byte{34, 92, 117, 48, 48, 48, 49, 92, 117, 48, 48, 48, 50, 92, 117, 48, 48, 48, 51, 92, 117, 48, 48, 48, 52, 34}

	result, err := MessageToPayload(input)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if !bytes.Equal(expected, result) {
		t.Errorf("Unexpected result: expected %v, got %v", expected, result)
	}
}

func TestInterfaceToBytes_Map(t *testing.T) {
	input := map[string]interface{}{
		"name": "John",
		"age":  30,
		"pets": []string{"dog", "cat"},
	}

	expected := []byte(`{
  "age": 30,
  "name": "John",
  "pets": [
    "dog",
    "cat"
  ]
}`)
	result, err := MessageToPayload(input)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if !bytes.Equal(expected, result) {
		t.Errorf("Unexpected result: expected %v, got %v", expected, result)
	}
}

func TestInterfaceToBytes_Slice(t *testing.T) {
	input := []interface{}{"a", 1, true}

	expected := []byte(`[
  "a",
  1,
  true
]`)

	result, err := MessageToPayload(input)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if !bytes.Equal(expected, result) {
		t.Errorf("Unexpected result: expected %v, got %v", expected, result)
	}
}

func TestInterfaceToBytes_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	input := Person{Name: "John", Age: 30}

	expected := []byte(`{
  "Name": "John",
  "Age": 30
}`)

	result, err := MessageToPayload(input)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if !bytes.Equal(expected, result) {
		t.Errorf("Unexpected result: expected %v, got %v", expected, result)
	}
}

package name

import (
	"testing"
)

func TestName_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{"v1.orders.process", true},
		{"v2.users.sync", true},
		{"v1.inventory.restock.queue", true},
		{"v1.orders.process.queue", true},
		{"v2.users.sync.queue", true},
		{"v1.orders.process_queue", false},
		{"orders.process", false},
		{"v1.orders", false},
		{"v1_orders.process", false},
	}

	for _, tc := range testCases {
		n := Name(tc.name)
		result := true

		if err := n.Validate(); err != nil {
			result = false
		}

		if result != tc.expected {
			t.Errorf("Expected Validate() for '%s' to be %v, got %v", tc.name, tc.expected, result)
		}
	}
}

func TestName_Parts(t *testing.T) {
	n := Name("v1.orders.process.queue")

	expected := []string{"v1.orders.process.queue", "v1", "orders", "process", "queue"}

	parts := n.Parts()

	for i, part := range parts {
		if part != expected[i] {
			t.Errorf("Expected part %d to be '%s', got '%s'", i, expected[i], part)
		}
	}
}

//nolint:gocritic
func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr bool
	}{
		{"v1.orders.process", false},
		{"v2.users.sync", false},
		{"v1.inventory.restock.queue", false},
		{"v1.orders.process.queue", false},
		{"v2.users.sync.queue", false},
		{"v1.orders.process_queue", true},
		{"orders.process", true},
		{"v1.orders", true},
		{"v1_orders.process", true},
	}

	for _, tc := range testCases {
		n, err := NewFromString(tc.name)

		if tc.expectedErr && err == nil {
			t.Errorf("Expected New() with name '%s' to return an error, got nil", tc.name)
		} else if !tc.expectedErr && err != nil {
			t.Errorf("Expected New() with name '%s' to return no error, got %v", tc.name, err)
		} else if !tc.expectedErr && n.String() != tc.name {
			t.Errorf("Expected New() to return Name with value '%s', got '%s'", tc.name, n.String())
		}
	}
}

func TestName_ToQueue(t *testing.T) {
	tests := []struct {
		name string
		n    Name
		want Queue
	}{
		{
			name: "name with .queue suffix",
			n:    Name("v1.orders.process.queue"),
			want: Queue("v1.orders.process.queue"),
		},
		{
			name: "name without .queue suffix",
			n:    Name("v1.orders.process"),
			want: Queue("v1.orders.process.queue"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ToQueue(); got != tt.want {
				t.Errorf("Name.ToQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName_ToTopic(t *testing.T) {
	tests := []struct {
		name string
		n    Name
		want Topic
	}{
		{
			name: "name with .queue suffix",
			n:    Name("v1.orders.process.queue"),
			want: Topic("v1.orders.process"),
		},
		{
			name: "name without .queue suffix",
			n:    Name("v1.orders.process"),
			want: Topic("v1.orders.process"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ToTopic(); got != tt.want {
				t.Errorf("Name.ToTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

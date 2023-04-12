package pubsub

import (
	"testing"
)

// This TEST exist just to ensure Mock match the IPubSub interface.
func TestMock_match_interface(t *testing.T) {
	var ips IPubSub = &Mock{}

	t.Log(ips)
}

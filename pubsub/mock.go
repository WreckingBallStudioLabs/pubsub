package pubsub

// Mock is a struct which satisfies the pubsub.IPubSub interface.
type Mock struct {
	// Publish sends a message to a topic.
	MockPublish func(topic string, message any) error

	// Subscribe subscribes to a topic and returns a channel for receiving messages.
	MockSubscribe func(topic, queue string, cb func([]byte)) Subscription

	// Unsubscribe unsubscribes from a topic.
	MockUnsubscribe func(topic string) error

	// Close closes the connection to the Pub Sub broker.
	MockClose func() error

	// GetClient returns the storage client. Use that to interact with the
	// underlying storage client.
	MockGetClient func() any
}

// Publish mocked call.
func (m *Mock) Publish(topic string, message any) error {
	return m.MockPublish(topic, message)
}

// Subscribe mocked call.
func (m *Mock) Subscribe(topic, queue string, cb func([]byte)) Subscription {
	return m.MockSubscribe(topic, queue, cb)
}

// Unsubscribe mocked call.
func (m *Mock) Unsubscribe(topic string) error {
	return m.MockUnsubscribe(topic)
}

// Close mocked call.
func (m *Mock) Close() error {
	return m.MockClose()
}

// GetClient mocked call.
func (m *Mock) GetClient() any {
	return m.MockGetClient()
}

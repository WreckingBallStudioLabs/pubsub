package pubsub

import "github.com/thalesfsp/sypl"

//////
// Creates the a struct which satisfies the storage.IStorage interface.
//////

// Mock is a struct which satisfies the pubsub.IPubSub interface.
type Mock struct {
	// Publish sends a message to a topic.
	MockPublish func(topic string, message any) error

	// PublishAsync sends a message to a topic. In case of error it will just log
	// it.
	MockPublishAsync func(topic string, message any)

	// Subscribe subscribes to a topic and returns a channel for receiving messages.
	MockSubscribe func(topic, queue string, cb func([]byte)) (Subscription, error)

	// Unsubscribe unsubscribes from a topic.
	MockUnsubscribe func(topic string) error

	// Close closes the connection to the Pub Sub broker.
	MockClose func() error

	// GetClient returns the storage client. Use that to interact with the
	// underlying storage client.
	MockGetClient func() any

	// GetLogger returns the logger.
	MockGetLogger func() sypl.ISypl

	// GetName returns the storage name.
	MockGetName func() string
}

//////
// When the methods are called, it will call the corresponding method in the
// Mock struct returning the desired value. This implements the IStorage
// interface.
//////

// Publish mocked call.
func (m *Mock) Publish(topic string, message any) error {
	return m.MockPublish(topic, message)
}

// PublishAsync sends a message to a topic. In case of error it will just log
// it.
func (m *Mock) PublishAsync(topic string, message any) {
	m.MockPublishAsync(topic, message)
}

// Subscribe mocked call.
func (m *Mock) Subscribe(topic, queue string, cb func([]byte)) (Subscription, error) {
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

// GetLogger returns the logger.
func (m *Mock) GetLogger() sypl.ISypl {
	return m.MockGetLogger()
}

// GetName returns the storage name.
func (m *Mock) GetName() string {
	return m.MockGetName()
}

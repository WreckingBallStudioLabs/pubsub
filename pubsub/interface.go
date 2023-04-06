package pubsub

import (
	"github.com/WreckingBallStudioLabs/pubsub/subscription"
	"github.com/thalesfsp/sypl"
)

//////
// Vars, consts, and types.
//////

// IPubSub defines a PubSub does.
type IPubSub interface {
	// Publish sends a message to a topic.
	Publish(topic string, message any) error

	// PublishAsync sends a message to a topic. In case of error it will just
	// log it.
	PublishAsync(topic string, message any)

	// Subscribe to a topic.
	Subscribe(topic, queue string, cb func([]byte)) (subscription.Subscription, error)

	// Unsubscribe from a topic.
	Unsubscribe(topic string) error

	// Close the connection to the Pub Sub broker.
	Close() error

	// GetClient returns the storage client. Use that to interact with the
	// underlying storage client.
	GetClient() any

	// GetLogger returns the logger.
	GetLogger() sypl.ISypl

	// GetName returns the pubsub name.
	GetName() string
}

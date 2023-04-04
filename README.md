# Pub Sub mechanism
The Pub Sub mechanism in our system provides a way for components to communicate with each other asynchronously. It consists of an interface IPubSub that defines the methods for interacting with the Pub Sub broker.

# Usage
To use the Pub Sub mechanism, you should first create an instance of the IPubSub interface. You can then use the following methods to interact with the Pub Sub broker:

`Publish(topic string, message any) error`

The Publish method sends a message to a topic. The topic parameter is a string that identifies the topic to which the message should be sent, and the message parameter is the message to be sent.

`Subscribe(topic, queue string, cb func([]byte)) Subscription`

The Subscribe method subscribes to a topic and returns a channel for receiving messages. The topic parameter is a string that identifies the topic to which the component wants to subscribe, and the queue parameter is a string that identifies the queue to which the component wants to subscribe. The cb parameter is a function that will be called whenever a message is received on the subscribed topic. The function takes a single parameter of type []byte, which is the message that was received. The method returns a Subscription object that can be used to unsubscribe from the topic.

`Unsubscribe(topic string) error`

The Unsubscribe method unsubscribes from a topic. The topic parameter is a string that identifies the topic from which the component wants to unsubscribe.

`Close() error`

The Close method closes the connection to the Pub Sub broker.

`GetClient() any`

The GetClient method returns the storage client. Use that to interact with the underlying storage client.

# Subscribing to Topics and Queues

To receive messages on a topic, a service should subscribe to that topic by providing a callback function that will be called whenever a message is received. This can be done using the Subscribe method of the IPubSub interface.

Example

Here's an example of how to subscribe to topics and queues using the nats Pub Sub implementation:
```go
func initSubscriptions() {
    // Get the Pub Sub client
    client := nats.Get()

    // Define the subscriptions
    subscriptions := []pubsub.Subscription{
        {
            Topic:    "onCampaignCreated",
            Queue:    "onCampaignCreatedQueue",
            Callback: OnCampaignCreated,
        },
        {
            Topic:    "onCampaignUpdated",
            Queue:    "onCampaignUpdatedQueue",
            Callback: OnCampaignUpdated,
        },
    }

    // Subscribe to each topic and queue
    for _, sub := range subscriptions {
        sub := client.Subscribe(sub.Topic, sub.Queue, sub.Callback)

        go func() {
            for msg := range sub.Channel {
                sub.Callback(msg)
            }
        }()
    }
}
```

This example subscribes to the "onCampaignCreated" and "onCampaignUpdated" topics using the OnCampaignCreated and OnCampaignUpdated callback functions, respectively. The Subscribe method returns a Subscription object, which is used to unsubscribe from the topic later if needed. The go statement creates a new goroutine that reads messages from the subscription's channel and passes them to the callback function.

To subscribe to topics and queues, simply call the initSubscriptions function from the resource's New() function or another appropriate location.

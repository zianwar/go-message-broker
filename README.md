## go-message-broker

A minimal implementation of a message broker in Go which satisfies the following interface:

```golang
type Message interface{}

type Broker interface {
	Subscribe(channel string) (<-chan Message, error)
	Unsubscribe(channel string) error
	Publish(channel string, m Message) error
	Close() error
}
```

### Imlementations:
- In-memory using go channels
- Redis PubSub
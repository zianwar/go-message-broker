package broker

type Message interface{}

type Subscriber interface {
	Subscribe(channel string) (<-chan Message, error)
	Unsubscribe(channel string) error
}

type Publisher interface {
	Publish(channel string, m Message) error
}

type Broker interface {
	Subscriber
	Publisher
	Close() error
}

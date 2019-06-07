package broker

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisBroker struct {
	*redis.Client
	subs map[string]*redis.PubSub
}

func NewRedisBroker() *redisBroker {
	client := redis.NewClient(&redis.Options{Addr: ":6379"})
	return &redisBroker{
		Client: client,
		subs:   make(map[string]*redis.PubSub),
	}
}

func (b *redisBroker) Subscribe(channel string) (<-chan Message, error) {
	sub := b.Client.Subscribe(channel)
	_, err := sub.Receive()
	if err != nil {
		return nil, err
	}

	b.subs[channel] = sub

	out := make(chan Message)
	go func() {
		for msg := range sub.Channel() {
			out <- msg
		}
		close(out)
	}()

	return out, nil
}

func (b *redisBroker) Unsubscribe(channel string) error {
	sub, ok := b.subs[channel]
	if !ok {
		return errors.Errorf("cannot find channel %s", channel)
	}

	if err := sub.Close(); err != nil {
		return err
	}

	return nil
}

func (b *redisBroker) Publish(channel string, m Message) error {
	_, ok := b.subs[channel]
	if !ok {
		return errors.Errorf("cannot find channel %s", channel)
	}

	if err := b.Client.Publish(channel, m).Err(); err != nil {
		return err
	}

	return nil
}

func (b *redisBroker) Close() error {
	if err := b.Client.Close(); err != nil {
		return err
	}
	return nil
}

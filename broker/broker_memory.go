package broker

import (
	"github.com/pkg/errors"
)

type memoryBroker struct {
	chans map[string]chan Message
}

func NewMemoryBroker() *memoryBroker {
	return &memoryBroker{
		chans: make(map[string]chan Message),
	}
}

func (b *memoryBroker) Subscribe(channel string) (<-chan Message, error) {
	ch := make(chan Message, 100)
	b.chans[channel] = ch
	return ch, nil
}

func (b *memoryBroker) Unsubscribe(channel string) error {
	ch, ok := b.chans[channel]
	if !ok {
		return errors.Errorf("cannot find channel %s", channel)
	}
	close(ch)
	delete(b.chans, channel)
	return nil
}

func (b *memoryBroker) Publish(channel string, m Message) error {
	ch, ok := b.chans[channel]
	if !ok {
		return errors.Errorf("cannot find channel %s", channel)
	}
	ch <- m
	return nil
}

func (b *memoryBroker) Close() error {
	for name, ch := range b.chans {
		close(ch)
		delete(b.chans, name)
	}
	return nil
}

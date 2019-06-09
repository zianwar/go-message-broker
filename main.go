package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zianwar/go-message-broker/broker"
)

type myMessage string

func main() {
	var b broker.Broker

	// Trying the in-memory broker.
	b = broker.NewMemoryBroker()
	// b = broker.NewRedisBroker()

	// subCh is a readony channel that we will
	// receive messages published on "ch1".
	subCh, err := b.Subscribe("ch1")
	if err != nil {
		log.Fatalln(err)
	}

	// start a publish loop
	// publish a message every second.
	go func() {
		defer b.Close()

		i := 0
		for {
			i++
			if err := b.Publish("ch1", fmt.Sprintf("message %d", i)); err != nil {
				log.Fatalln(err)
			}

			time.Sleep(time.Second)
			// stop after 5 iterations.
			if i == 5 {
				if err := b.Unsubscribe("ch1"); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	}()

	// read messages from subCh published on "ch1".
	for m := range subCh {
		fmt.Printf("got message: %s\n", m)
	}

}

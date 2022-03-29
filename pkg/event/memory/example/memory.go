package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Yugive/goat_kit/pkg/event"
	"github.com/Yugive/goat_kit/pkg/event/memory"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	producer, consumer := memory.NewMemory("test")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	errCh := make(chan error, 1)

	go receive(consumer, errCh)

loop:
	for i := 0; i < 5; i++ {
		msg := memory.NewMessage("test", []byte("hello world"))
		err := producer.Send(context.Background(), msg)
		if err != nil {
			errCh <- errors.New("producer error")
			break loop
		}
	}

	select {
	case <-errCh:
		fmt.Println("consume msg error")
	case <-sigs:
		fmt.Println("exit")
	}
	_ = producer.Close()
	_ = consumer.Close()
}

func receive(c event.Consumer, errCh chan error) {
	fmt.Println("Starting...")
	err := c.Receive(context.Background(), func(ctx context.Context, msg event.Event) error {
		fmt.Printf("k:%s, v:%s\n", msg.Key(), msg.Value())
		return nil
	})

	if err != nil {
		errCh <- err
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/Yugive/goat_kit/pkg/event"
	"github.com/Yugive/goat_kit/pkg/event/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	receiver, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "kira", kafka.GroupID("group-test"))
	if err != nil {
		log.Fatal(err)
	}
	errCh := make(chan error, 1)
	go receive(receiver, errCh)
	select {
	case <-errCh:
		fmt.Println("consume msg error")
	case <-sigs:
		fmt.Println("exit")
	}
	_ = receiver.Close()
}

func receive(receiver event.Consumer, errCh chan error) {
	fmt.Println("Starting...")
	err := receiver.Receive(context.Background(), func(ctx context.Context, msg event.Event) error {
		fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
		return nil
	})
	if err != nil {
		errCh <- err
	}
}

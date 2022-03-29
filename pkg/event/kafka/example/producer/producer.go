package main

import (
	"context"
	"fmt"
	"github.com/Yugive/goat_kit/pkg/event"
	"github.com/Yugive/goat_kit/pkg/event/kafka"
	"log"
)

func main() {
	sender, err := kafka.NewKafkaWriter([]string{"localhost:9092"}, "kira")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		send(sender)
	}

	_ = sender.Close()
}

func send(sender event.Producer) {
	msg := kafka.NewMessage("lope", []byte("hello world"))
	err := sender.Send(context.Background(), msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
}

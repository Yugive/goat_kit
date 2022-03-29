package event

import "context"

type Event interface {
	Key() string
	Value() []byte
}

type Handler func(context.Context, Event) error

type Producer interface {
	Send(ctx context.Context, msg Event) error
	Close() error
}

type Consumer interface {
	Receive(ctx context.Context, handler Handler) error
	Close() error
}

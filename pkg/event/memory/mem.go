package memory

import (
	"context"
	"github.com/ngyugive/goat_kit/pkg/event"
	"sync"
)

var (
	_ event.Producer = (*memoryProducer)(nil)
	_ event.Consumer = (*memoryReceiver)(nil)
	_ event.Event    = (*Message)(nil)

	chanMap = struct {
		sync.RWMutex
		chm map[string]chan *Message
	}{}

	chanSize = 1024
)

func init() {
	chanMap.chm = make(map[string]chan *Message)
}

type Message struct {
	key   string
	value []byte
}

func (m *Message) Key() string {
	return m.key
}

func (m *Message) Value() []byte {
	return m.value
}

func NewMessage(key string, value []byte) event.Event {
	return &Message{
		key:   key,
		value: value,
	}
}

type memoryProducer struct {
	topic string
}

func (m *memoryProducer) Send(ctx context.Context, msg event.Event) error {
	chanMap.chm[m.topic] <- &Message{
		key:   msg.Key(),
		value: msg.Value(),
	}
	return nil
}

func (m *memoryProducer) Close() error {
	return nil
}

type memoryReceiver struct {
	topic string
}

func (m *memoryReceiver) Receive(ctx context.Context, handler event.Handler) error {
	for msg := range chanMap.chm[m.topic] {
		err := handler(context.Background(), msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *memoryReceiver) Close() error {
	return nil
}

func NewMemory(topic string) (event.Producer, event.Consumer) {
	chanMap.Lock()
	defer chanMap.Unlock()
	if _, ok := chanMap.chm[topic]; !ok {
		chanMap.chm[topic] = make(chan *Message, chanSize)
	}
	return &memoryProducer{topic: topic}, &memoryReceiver{topic: topic}
}

package kafka

import (
	"context"
	"crypto/tls"
	"github.com/Yugive/goat_kit/pkg/event"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"time"
)

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

type WriterOption func(*kafka.Writer)

func MaxAttempts(times int) WriterOption {
	return func(p *kafka.Writer) {
		p.MaxAttempts = times
	}
}

func BatchSize(size int) WriterOption {
	return func(p *kafka.Writer) {
		p.BatchSize = size
	}
}

func BatchBytes(size int64) WriterOption {
	return func(p *kafka.Writer) {
		p.BatchBytes = size
	}
}

func BatchTimeout(timeout time.Duration) WriterOption {
	return func(p *kafka.Writer) {
		p.BatchTimeout = timeout
	}
}

func ReadTimeout(timeout time.Duration) WriterOption {
	return func(p *kafka.Writer) {
		p.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) WriterOption {
	return func(p *kafka.Writer) {
		p.WriteTimeout = timeout
	}
}

func Async(need bool) WriterOption {
	return func(p *kafka.Writer) {
		p.Async = need
	}
}

type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func (p *kafkaProducer) Send(ctx context.Context, message event.Event) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.Key()),
		Value: message.Value(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *kafkaProducer) Close() error {
	err := p.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaProducer(addr []string, topic string, opts ...WriterOption) (event.Producer, error) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(addr...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	for _, opt := range opts {
		opt(w)
	}

	return &kafkaProducer{writer: w, topic: topic}, nil
}

type WriterConfig func(*kafka.WriterConfig)

func MaxWriterAttempts(times int) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.MaxAttempts = times
	}
}

func WriterBatchSize(size int) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.BatchSize = size
	}
}

func WriterBatchBytes(size int) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.BatchBytes = size
	}
}

func WriterBatchTimeout(timeout time.Duration) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.BatchTimeout = timeout
	}
}

func WriterWriteTimeout(timeout time.Duration) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.WriteTimeout = timeout
	}
}

func WriterReadTimeout(timeout time.Duration) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.ReadTimeout = timeout
	}
}

func WriterAsync(need bool) WriterConfig {
	return func(c *kafka.WriterConfig) {
		c.Async = need
	}
}

func TLS(conf *tls.Config) WriterConfig {
	return func(c *kafka.WriterConfig) {
		if c.Dialer != nil {
			c.Dialer.TLS = conf
		} else {
			c.Dialer = &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
				TLS:       conf,
			}
		}
	}
}

func SASL(username string, password string) WriterConfig {
	return func(c *kafka.WriterConfig) {
		mechanism := plain.Mechanism{
			Username: username,
			Password: password,
		}
		if c.Dialer != nil {
			c.Dialer.SASLMechanism = mechanism
		} else {
			c.Dialer = &kafka.Dialer{
				Timeout:       10 * time.Second,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
		}
	}
}

func NewKafkaWriter(addr []string, topic string, opts ...WriterConfig) (event.Producer, error) {
	conf := kafka.WriterConfig{
		Brokers:  addr,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	for _, opt := range opts {
		opt(&conf)
	}

	w := kafka.NewWriter(conf)
	return &kafkaProducer{writer: w, topic: topic}, nil
}

type kafkaConsumer struct {
	reader *kafka.Reader
	topic  string
}

func (c *kafkaConsumer) Receive(ctx context.Context, handler event.Handler) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		err = handler(ctx, &Message{
			key:   string(msg.Key),
			value: msg.Value,
		})
		if err != nil {
			return err
		}
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func (c *kafkaConsumer) Close() error {
	err := c.reader.Close()
	if err != nil {
		return err
	}
	return nil
}

type ReaderOption func(reader *kafka.ReaderConfig)

func GroupID(id string) ReaderOption {
	return func(r *kafka.ReaderConfig) {
		r.GroupID = id
	}
}

func MinBytes(size int) ReaderOption {
	return func(r *kafka.ReaderConfig) {
		r.MinBytes = size
	}
}

func MaxBytes(size int) ReaderOption {
	return func(r *kafka.ReaderConfig) {
		r.MaxBytes = size
	}
}

func StartOffset(offset int64) ReaderOption {
	return func(r *kafka.ReaderConfig) {
		r.StartOffset = offset
	}
}

func ReaderMaxAttempts(times int) ReaderOption {
	return func(r *kafka.ReaderConfig) {
		r.MaxAttempts = times
	}
}

func ReaderPartition(p int) ReaderOption {
	return func(r *kafka.ReaderConfig) {
		r.Partition = p
	}
}

func ReaderTLS(conf *tls.Config) ReaderOption {
	return func(c *kafka.ReaderConfig) {
		if c.Dialer != nil {
			c.Dialer.TLS = conf
		} else {
			c.Dialer = &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
				TLS:       conf,
			}
		}
	}
}

func ReaderSASL(username string, password string) ReaderOption {
	return func(c *kafka.ReaderConfig) {
		mechanism := plain.Mechanism{
			Username: username,
			Password: password,
		}
		if c.Dialer != nil {
			c.Dialer.SASLMechanism = mechanism
		} else {
			c.Dialer = &kafka.Dialer{
				Timeout:       10 * time.Second,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
		}
	}
}

func NewKafkaConsumer(addr []string, topic string, opts ...ReaderOption) (event.Consumer, error) {
	c := kafka.ReaderConfig{
		Brokers: addr,
		Topic:   topic,
	}
	for _, opt := range opts {
		opt(&c)
	}
	r := kafka.NewReader(c)
	return &kafkaConsumer{reader: r, topic: topic}, nil
}

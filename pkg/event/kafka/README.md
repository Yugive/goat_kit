# Kafka SDK

# 依赖库
* segmentio/kafka-go [github.com/segmentio/kafka-go]

# 调用示例
* producer
  
  producer, _ := kafka.NewKafkaWriter([]string{"localhost:9092"}, "test")
  
  msg := kafka.NewMessage("test_key", []byte("hello world"))

  producer.Send(context.Background(), msg)

* consumer

  receiver, _ := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "test", kafka.GroupID("group-test"))

  err := receiver.Receive(context.Background(), func(ctx context.Context, msg event.Event) error {
  
  fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
  
  return nil
  
  })
  
  
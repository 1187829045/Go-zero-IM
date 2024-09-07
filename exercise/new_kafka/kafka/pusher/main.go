package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	// 配置Kafka生产者
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// 创建生产者
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	defer producer.Close()

	// 构建消息
	msg := &sarama.ProducerMessage{
		Topic: "example_topic",
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	// 发送消息
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	} else {
		log.Printf("Message sent to partition %d with offset %d\n", partition, offset)
	}
}

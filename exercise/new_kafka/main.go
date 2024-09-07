package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	// 创建一个 Kafka 生产者
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",       // Kafka broker 地址
		"key.serializer":    kafka.StringSerializer, // 键的序列化方式
		"value.serializer":  kafka.StringSerializer, // 值的序列化方式
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// 要发送的消息
	topic := "example_topic"
	key := "example_key"
	value := "Hello, Kafka!"

	// 创建一个 Kafka 消息
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
	}

	// 发送消息
	err = producer.Produce(msg, nil)
	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	}

	// 等待消息被成功送达
	producer.Flush(15 * 1000) // 15秒

	fmt.Println("Message sent successfully!")
}

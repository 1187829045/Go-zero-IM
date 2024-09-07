package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	// 配置Kafka消费者
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama consumer: %v", err)
	}
	defer consumer.Close()

	// 订阅主题
	partitionConsumer, err := consumer.ConsumePartition("example_topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// 消费消息
	for msg := range partitionConsumer.Messages() {
		log.Printf("Received message: %s from partition %d at offset %d\n", string(msg.Value), msg.Partition, msg.Offset)
	}
}

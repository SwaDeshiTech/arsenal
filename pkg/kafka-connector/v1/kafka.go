package v1

import (
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConfig struct {
	URL     string `yaml:"url"`
	GroupID string `yaml:"groupID"`
}

type MessageHandler struct {
	Producer *kafka.Producer `json:"producer"`
	Topic    string          `json:"topic"`
	Message  string          `json:"message"`
}

func (kafkaConfig *KafkaConfig) InitProducer() (*kafka.Producer, error) {

	var lock sync.Mutex

	lock.Lock()
	defer lock.Unlock()

	config := kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.URL,
		"group.id":          kafkaConfig.GroupID,
	}

	p, err := kafka.NewProducer(&config)
	if err != nil {
		log.Println("failed to connect to kafka", err)
		return nil, err
	}

	return p, nil
}

func (messageHandler *MessageHandler) SendMessageHandler() {

	go func() {
		for e := range messageHandler.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Println("message delivered to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// produce messages to topic (asynchronously)
	err := messageHandler.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &messageHandler.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(messageHandler.Message),
	}, nil)
	if err != nil {
		log.Println("failed to produce message on kafka", err)
		messageHandler.Producer.Flush(0)
		return
	}
	// wait for message deliveries before shutting down
}

func (kafkaConfig *KafkaConfig) GetConsumer() (*kafka.Consumer, error) {

	config := kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.URL,
		"group.id":          kafkaConfig.GroupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Println("failed to get kafka consumer", err)
		return nil, err
	}
	return consumer, nil
}

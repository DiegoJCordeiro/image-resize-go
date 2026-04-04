package kafka

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ImageKafkaProducer struct {
	producer *kafka.Producer
}

func NewImageProducerImpl(bootstrapServer string) *ImageKafkaProducer {

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":   bootstrapServer,
		"delivery.timeout.ms": 6000,
		"acks":                "all",
		"enable.idempotence":  "true",
	}

	kafkaProducer, err := kafka.NewProducer(kafkaConfig)

	if err != nil {
		panic(err)
	}

	return &ImageKafkaProducer{
		producer: kafkaProducer,
	}
}

func (ep *ImageKafkaProducer) Publish(jsonEvent, topic string, key []byte) error {

	eventChan := make(chan kafka.Event, 5)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:       key,
		Value:     []byte(jsonEvent),
		Timestamp: time.Now(),
	}

	err := ep.producer.Produce(msg, eventChan)

	go func() {
		ep.deliveryReport(eventChan)
	}()

	if err != nil {
		return err
	}

	return nil
}

func (ep *ImageKafkaProducer) deliveryReport(eventReport <-chan kafka.Event) error {

	reportReceived := <-eventReport
	_, ok := reportReceived.(*kafka.Message)

	if !ok {
		return fmt.Errorf("ReportReceived: not message")
	}

	return nil
}

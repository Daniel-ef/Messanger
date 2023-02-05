package kafkacontroller

import (
	"context"
	"fmt"
	"github.com/messanger/services/messaging/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type Topic string

const (
	TopicMessages Topic = "messages"
)

type KafkaController struct {
	writers map[Topic]*kafka.Writer
	readers map[Topic]*kafka.Reader
}

func NewKafkaController(ctx context.Context, addresses ...string) (*KafkaController, func(ctx context.Context)) {
	controller := &KafkaController{
		writers: make(map[Topic]*kafka.Writer),
		readers: make(map[Topic]*kafka.Reader),
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(addresses...),
		Topic:                  string(TopicMessages),
		AllowAutoTopicCreation: true,
	}
	controller.writers[TopicMessages] = w

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   addresses,
		Topic:     string(TopicMessages),
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	if err := r.SetOffsetAt(ctx, time.Now()); err != nil {
		logger.Error(ctx, "could not set offset at", zap.Error(err))
	}
	controller.readers[TopicMessages] = r

	closer := func(ctx context.Context) {
		err := w.Close()
		if err != nil {
			logger.Error(ctx, "could not close kafka writer", zap.Error(err))
		}
		err = r.Close()
		if err != nil {
			logger.Error(ctx, "could not close kafka reader", zap.Error(err))
		}
	}

	return controller, closer
}

func (k *KafkaController) SendMessage(ctx context.Context, topic Topic, msg []byte) error {
	err := k.writers[topic].WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		return fmt.Errorf("could not write message to kafka: %w", err)
	}

	return nil
}

func (k *KafkaController) ReadMessage(ctx context.Context, topic Topic) (kafka.Message, error) {
	message, err := k.readers[topic].ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("could not read message from kafka: %w", err)
	}

	return message, nil
}

func (k *KafkaController) CommitMessage(ctx context.Context, message kafka.Message) error {
	err := k.readers[Topic(message.Topic)].CommitMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("could not commit messages: %w", err)
	}

	return nil
}

package kafka

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	protos "github.com/jparkkennaby/docker-kafka-go/gen/go/protos"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
	envelope "pkg.dsb.dev/event/v1"
)

// Config defines the necessary configuration for the Kafka consumer.
type ConsumerConfig struct {
	Brokers       []string
	ConsumerGroup string
	Topic         string
	UseTLS        bool
}

// Consumer wraps the franz-go Kafka client.
type Consumer struct {
	client *kgo.Client
}

// NewConsumer initializes and returns a Franz-go Kafka consumer.
func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	clientOpts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.AutoCommitMarks(), // Enables auto-commit
	}

	client, err := kgo.NewClient(clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{client: client}, nil
}

type ConsumeMessageHandler func(context.Context, *protos.Event) error

func (c *Consumer) Consume(ctx context.Context, handler ConsumeMessageHandler) error {
	return c.HandleConsume(ctx, transformHandler(handler))
}

type RecordHandler func(context.Context, *kgo.Record) error

// Consume listens for messages and processes them.
func (c *Consumer) HandleConsume(ctx context.Context, handler RecordHandler) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			log.Println("Consumer closed.")
			return nil
		}

		errs := fetches.Errors()
		for _, err := range errs {
			log.Printf("Error consuming from topic %s: %v\n", err.Topic, err.Err)
		}

		fetches.EachRecord(func(r *kgo.Record) {
			log.Printf("Received message: %s from topic %s\n", string(r.Value), r.Topic)
			// Manually mark message as processed if auto-commit is disabled.
			c.client.MarkCommitRecords(r)
		})
	}
}

// Close shuts down the Kafka consumer.
func (c *Consumer) Close() {
	c.client.Close()
}

func transformHandler(handler ConsumeMessageHandler) func(context.Context, *kgo.Record) error {
	return func(ctx context.Context, record *kgo.Record) error {
		var env envelope.Envelope
		if err := proto.Unmarshal(record.Value, &env); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal event envelope", slog.Any("error", err))
			eventConsumedResult.WithLabelValues("unmarshal_error_1", "failed").Inc()
			return fmt.Errorf("failed to unmarshal event envelope: %w", err)
		}

		payload, err := env.Payload.UnmarshalNew()
		if err != nil {
			eventConsumedResult.WithLabelValues(env.GetPayload().GetTypeUrl(), "failed").Inc()
			if isPayloadTypeUnrecognizedError(err) {
				return nil
			}
			slog.ErrorContext(ctx, "failed to unmarshal event payload", slog.Any("error", err))
			return fmt.Errorf("failed to unmarshal event payload (event ID: %s): %w", env.Id, err)
		}

		eventConsumedResult.WithLabelValues(env.GetPayload().GetTypeUrl(), "success").Inc()

		//nolint: gocritic
		switch event := payload.(type) {
		case *protos.Event:
			return handler(ctx, event)
		}

		return nil
	}
}

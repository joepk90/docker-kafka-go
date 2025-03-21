package system

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jparkkennaby/docker-kafka-go/eventhandler"
	"github.com/jparkkennaby/docker-kafka-go/store"
	"github.com/jparkkennaby/docker-kafka-go/system/kafka"
	"golang.org/x/sync/errgroup"
)

type System struct {
	config *Config
}

func Initialize() (*System, error) {
	config, err := configFromEnv()
	if err != nil {
		usage()
		return nil, err
	}

	return &System{config: config}, nil
}

func (s *System) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return ServeMetricsAndHealthChecks(
			ctx,
			s.config.AppName,
			s.config.OperationalAddress,
		)
	})

	eg.Go(func() error {
		return runBConsumer(ctx, s.config)
	})

	return eg.Wait()
}

func runBConsumer(ctx context.Context, config *Config) error {
	producer, err := kafka.NewProducer(kafka.KafkaProducerConfig(config.KafkaConfig.Producer))
	if err != nil {
		return fmt.Errorf("kafka producer: %w", err)
	}
	defer producer.Close()

	dlqProducer, err := kafka.NewProducer(kafka.KafkaProducerConfig(config.KafkaConfig.DLQProducer))
	if err != nil {
		return fmt.Errorf("kafka dlq producer: %w", err)
	}
	defer dlqProducer.Close()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	db, err := store.NewDB(timeoutCtx, store.Config{DatabaseURL: config.DatabaseConnectionString})
	if err != nil {
		return fmt.Errorf("setup cockroach: %v", err)
	}

	handler := &eventhandler.Handler{
		Producer: producer,
		Store:    db,
	}

	consumeHandler := kafka.ConsumeMessageHandler(
		eventhandler.Wrap(
			handler.OnEvent,
			eventhandler.Recover,
			eventhandler.Tracing,
			eventhandler.ErrorHandling(dlqProducer),
		))

	kafkaConsumerConfig := kafka.ConsumerConfig(config.KafkaConfig.Consumer)
	kafkaConsumer, err := kafka.NewConsumer(kafkaConsumerConfig)
	if err != nil {
		return fmt.Errorf("kafka consumer: %w", err)
	}
	defer kafkaConsumer.Close()

	return kafkaConsumer.Consume(ctx, consumeHandler)
}

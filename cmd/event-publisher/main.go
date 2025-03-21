package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jparkkennaby/docker-kafka-go/cmd/event-publisher/config"
	"github.com/jparkkennaby/docker-kafka-go/system/kafka"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.FromEnv()
	if err != nil {
		config.Usage()
		return err
	}

	kakfaProducer, err := kafka.NewProducer(kafka.KafkaProducerConfig(cfg.Producer))
	if err != nil {
		return fmt.Errorf("bill delivered event producer: %w", err)
	}
	defer kakfaProducer.Close()

	handler := NewHandler(kakfaProducer)

	err = handler.Handle(ctx)

	if err != nil {
		return fmt.Errorf("something went wrong: %w", err)
	}

	return nil
}

package main

import (
	"context"
	"fmt"

	protos "github.com/jparkkennaby/docker-kafka-go/gen/go/protos"
	"github.com/jparkkennaby/docker-kafka-go/system/kafka"
)

const (
	AppName = "docker-kafka-go-event-publisher"
)

type Handler struct {
	producer *kafka.Producer
}

func NewHandler(producer *kafka.Producer) *Handler {
	return &Handler{
		producer: producer,
	}
}

func (h *Handler) Handle(ctx context.Context) error {
	deliveryRequest := &protos.Event{
		Id:     "123",
		Tag:    protos.Tag_EVENT_TAG_TECHNOLOGY,
		Status: protos.Status_EVENT_STATUS_SUCCESSFUL,
		// CreatedAt: *timestamppb.Timestamp,
	}

	err := h.producer.ProduceMessage(ctx, deliveryRequest, AppName, "protos.Event")
	if err != nil {
		return fmt.Errorf("failed to send to kafka: %w", err)
	}

	return nil
}

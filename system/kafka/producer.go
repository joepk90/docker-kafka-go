package kafka

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jparkkennaby/docker-kafka-go/system/tracing"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	envelope "pkg.dsb.dev/event/v1"
)

var tracer = noop.NewTracerProvider().Tracer("")

type KafkaProducerConfig struct {
	Brokers []string
	Topic   string
	UseTLS  bool
}

type Producer struct {
	client *kgo.Client
}

func NewProducer(cfg KafkaProducerConfig) (*Producer, error) {
	clientOpts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.DefaultProduceTopic(cfg.Topic),
		// kgo.ConsumerGroup(cfg.ConsumerGroup), // TODO REVIEW LINE
		// kgo.AutoCommitMarks(), // Enables auto-commit // TODO REVIEW LINE
	}

	prioOpts := append(clientOpts, defaultNativeOpts()...)

	client, err := kgo.NewClient(prioOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed creating franz-go client: %w", err)
	}

	return &Producer{
		client: client,
	}, nil
}

func (p *Producer) ProduceMessage(ctx context.Context, message proto.Message, app, comment string) (err error) {
	ctx, span := tracer.Start(ctx, "produce-message")
	defer span.End()

	span.AddEvent("marshal-envelope-payload")
	envelopePayload, err := anypb.New(message)
	if err != nil {
		return tracing.RecordAndReturnError(span, fmt.Errorf("failed to marshal event envelopePayload: %w", err))
	}

	metadata := map[string]string{
		"domain": "docker-kafka-go",
	}

	envelopedEvent := &envelope.Envelope{
		Id:        uuid.NewString(),
		Timestamp: timestamppb.Now(),
		Payload:   envelopePayload,
		Sender: &envelope.Sender{
			Metadata:    metadata,
			Application: app,
		},
	}

	err = p.Produce(ctx, envelopedEvent)
	if err != nil {
		return tracing.RecordAndReturnError(span, fmt.Errorf("failed to produce event: %w", err))
	}

	return nil
}

func (p *Producer) Produce(ctx context.Context, event *envelope.Envelope) error {
	ctx, span := tracer.Start(ctx, "produce-enveloped-event")
	defer span.End()

	span.AddEvent("marshal-enveloped-event")
	payload, err := proto.Marshal(event)
	if err != nil {
		return tracing.RecordAndReturnError(span, fmt.Errorf("failed to marshal message: %w", err))
	}

	rec := &kgo.Record{
		Context: ctx,
		Key:     []byte(event.Id),
		Value:   payload,
	}

	span.AddEvent("producing-event")
	res := p.client.ProduceSync(ctx, rec)
	if err = res.FirstErr(); err != nil {
		return tracing.RecordAndReturnError(span, fmt.Errorf("failed to produce message. first error: %w", err))
	}

	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}

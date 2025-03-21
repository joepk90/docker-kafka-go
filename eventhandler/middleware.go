package eventhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	protos "github.com/jparkkennaby/docker-kafka-go/gen/go/protos"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OnEvent func(context.Context, *protos.Event) error

var tracer trace.Tracer

func RecordAndReturnError(span trace.Span, err error) error {
	RecordError(span, err)
	return err
}

// RecordError records the error and marks the span as erroring.
// When the error is nil this is no op.
func RecordError(span trace.Span, err error) {
	// if you try and record nothing then we record nothing
	if err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func Tracing(f OnEvent) OnEvent {
	return func(ctx context.Context, event *protos.Event) error {
		ctx, span := tracer.Start(ctx, "event-handling")
		defer span.End()

		span.SetAttributes(attribute.KeyValue{
			Key:   "event_id",
			Value: attribute.StringValue(event.GetId()),
		},
			attribute.KeyValue{
				Key:   "event_status",
				Value: attribute.StringValue(event.GetStatus().String()),
			},
		)

		return RecordAndReturnError(span, f(ctx, event))
	}
}

func Recover(f OnEvent) OnEvent {
	return func(ctx context.Context, event *protos.Event) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Join(err, fmt.Errorf("recovered: %v", r))
			}
		}()
		return f(ctx, event)
	}
}

func ErrorHandling(dlq Producer) func(OnEvent) OnEvent {
	return func(f OnEvent) OnEvent {
		return func(ctx context.Context, event *protos.Event) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			start := time.Now()
			if err := f(ctx, event); err != nil {
				slog.ErrorContext(ctx, "event failed",
					slog.String("event_id", event.GetId()),
					slog.String("event_status", event.GetStatus().String()),
					slog.Duration("time", time.Since(start)),
					slog.Any("error", err),
				)
				// this would crash our service on errors!
				return dlq.ProduceMessage(ctx, event, AppName, err.Error())
			}

			return nil
		}
	}
}

type Middleware func(OnEvent) OnEvent

// Wrap the function f with the given middlewares.
// The order is important!
// Calling Wrap(f, g, h)(args) will execute:
// h(g(f(args)))
func Wrap(f OnEvent, mws ...Middleware) OnEvent {
	for _, mw := range mws {
		f = mw(f)
	}
	return f
}

package tracing

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

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

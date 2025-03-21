package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	eventConsumedResult = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "example-code",
			Subsystem: "docker-kafka-go",
			Name:      "event_consumed_result",
			Help:      "Counter that records the number of events consumed",
		},
		[]string{
			"event_type", "result",
		},
	)
)

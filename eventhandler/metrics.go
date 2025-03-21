package eventhandler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	eventHandledResult = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "example-code",
			Subsystem: "docker-kafka-go",
			Name:      "event_handled",
			Help:      "Records the result of event handler",
		},
		[]string{
			"result", "reason", "delivery_tag",
		},
	)
)

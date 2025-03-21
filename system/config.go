package system

import "github.com/kelseyhightower/envconfig"

const appName = "docker-kafka-go"

type Config struct {
	// OperationalAddress serves a health and a ready check for the probes in k8s
	// and a metrics endpoints that returns prometheus metrics.
	OperationalAddress string `envconfig:"OPERATIONAL_ADDR" default:"0.0.0.0:8081" required:"true"`

	// KafkaConfig is used to configure all kafka urls.
	KafkaConfig struct {
		Consumer struct {
			Brokers       []string `envconfig:"KAFKA_CONSUMER_BROKERS"`
			ConsumerGroup string   `envconfig:"KAFKA_CONSUMER_GROUP"`
			Topic         string   `envconfig:"KAFKA_CONSUMER_TOPIC"`
			// UseTLS        bool     `envconfig:"KAFKA_CONSUMER_USE_TLS" default:"false"`
		}

		Producer struct {
			Brokers []string `envconfig:"KAFKA_PRODUCER_BROKERS"`
			Topic   string   `envconfig:"KAFKA_PRODUCER_TOPIC"`
			// UseTLS  bool     `envconfig:"KAFKA_PRODUCER_USE_TLS" default:"false"`
		}

		DLQProducer struct {
			Brokers []string `envconfig:"KAFKA_DLQ_BROKERS"`
			Topic   string   `envconfig:"KAFKA_DLQ_TOPIC"`
			// UseTLS  bool     `envconfig:"KAFKA_DLQ_USE_TLS" default:"true"`
		}
	}

	DatabaseConnectionString string `envconfig:"DATABASE_CONNECTION_STRING" required:"true"`

	// Values that are not set by the environment:
	AppName string `ignored:"true"`
}

// ConfigFromEnv reads the configuration from the environment.
func configFromEnv() (*Config, error) {
	config := Config{}
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}
	config.AppName = appName
	return &config, nil
}

// Usage prints the configuration to stdout.
func usage() {
	config := Config{}
	_ = envconfig.Usage("", &config)
}

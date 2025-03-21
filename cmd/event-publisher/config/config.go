package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config contains all configuration for the invoice-fulfillment service.
type Config struct {
	AppName string `envconfig:"APP_NAME"`

	// defaults are set to work with local kafka instance / docker-compose setup
	Producer struct {
		Brokers []string `envconfig:"KAFKA_PRODUCER_BROKERS_EXT" default:"localhost:9094"`
		Topic   string   `envconfig:"KAFKA_PRODUCER_TOPIC_EXT" default:"event-topic-in"`
		// UseTLS  bool     `envconfig:"KAFKA_PRODUCER_USE_TLS" default:"false"`
	}
}

// FromEnv reads the configuration from the environment.
func FromEnv() (*Config, error) {
	config := Config{}
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Usage prints the configuration to stdout.
func Usage() {
	config := Config{}
	_ = envconfig.Usage("", &config)
}

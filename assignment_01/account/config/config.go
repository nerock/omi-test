package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	PostgresURI   string `json:"postgres_uri"`
	GrpcPort      int    `json:"grpc_port"`
	HttpPort      int    `json:"http_port"`
	LogLevel      int    `json:"log_level"`
	LogStackTrace bool   `json:"log_stacktrace"`
	NatsUri       string `json:"nats_uri"`
	NatsTopic     string `json:"nats_topic"`
}

func FromFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not decode config json: %w", err)
	}

	return cfg, nil
}

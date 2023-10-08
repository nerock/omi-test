package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Postgres PostgresConfig `json:"postgres"`
	Server   ServerConfig   `json:"server"`
	Log      LogConfig      `json:"log"`
	Nats     NatsConfig     `json:"nats"`
}

type PostgresConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

type ServerConfig struct {
	GrpcPort int `json:"grpc_port"`
	HttpPort int `json:"http_port"`
}

type LogConfig struct {
	Level      int  `json:"level"`
	StackTrace bool `json:"stacktrace"`
}

type NatsConfig struct {
	Uri   string `json:"uri"`
	Topic string `json:"topic"`
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

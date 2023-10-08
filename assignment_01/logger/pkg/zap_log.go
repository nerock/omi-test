package pkg

import (
	"github.com/nerock/omi-test/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(cfg config.LogConfig) (*zap.Logger, error) {
	return zap.Config{
		Level: zap.NewAtomicLevelAt(zapcore.Level(cfg.Level)),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: !cfg.StackTrace,
	}.Build()
}

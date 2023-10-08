package pkg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(lvl int, stackTrace bool) (*zap.Logger, error) {
	return zap.Config{
		Level: zap.NewAtomicLevelAt(zapcore.Level(lvl)),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: !stackTrace,
	}.Build()
}

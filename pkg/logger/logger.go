package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New(opts ...Option) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	return &Logger{log}, nil
}

package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
	fields []Field
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
	return &Logger{log, nil}, nil
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, append(l.fields, fields...)...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, append(l.fields, fields...)...)
}

func (l *Logger) WithFields(fields ...Field) *Logger {
	if len(fields) == 0 {
		return l
	}
	return &Logger{
		Logger: l.Logger,
		fields: append(l.fields, fields...),
	}
}

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config = zap.Config

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(cfg *Config) {
	f(cfg)
}

func WithTimeEncoderOfLayout(layout string) Option {
	return optionFunc(func(cfg *Config) {
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(layout)
	})
}

func WithRFC3339TimeEncoder() Option {
	return optionFunc(func(cfg *Config) {
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	})
}

func WithOutput(paths ...string) Option {
	return optionFunc(func(cfg *Config) {
		cfg.OutputPaths = paths
	})
}

func WithErrorOutput(paths ...string) Option {
	return optionFunc(func(cfg *Config) {
		cfg.ErrorOutputPaths = paths
	})
}

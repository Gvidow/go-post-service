package logger

import "go.uber.org/zap/zapcore"

type Field = zapcore.Field

func String(field, val string) Field {
	return zapcore.Field{
		Key:    field,
		String: val,
		Type:   zapcore.StringType,
	}
}

func Int(field string, val int) Field {
	return zapcore.Field{
		Key:     field,
		Integer: int64(val),
		Type:    zapcore.Int64Type,
	}
}

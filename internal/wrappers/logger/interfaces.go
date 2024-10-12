package logger

import "log/slog"

type Logger interface {
	SetLogLevel(level slog.Level)
	ParseLogLevel(level string) (slog.Level, error)
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
}

package logger

import (
	"fmt"
	"log/slog"
	"strings"
)

type SlogWrapper struct {
	logger   *slog.Logger
	logLevel *slog.LevelVar
}

// NewSlogWrapper creates a new instance of SlogWrapper with the given slog.Logger.
func NewSlogWrapper(logger *slog.Logger, logLevel *slog.LevelVar) *SlogWrapper {
	return &SlogWrapper{logger: logger, logLevel: logLevel}
}

func (l *SlogWrapper) SetLogLevel(level slog.Level) {
	l.logLevel.Set(level)
}

func (l *SlogWrapper) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *SlogWrapper) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, keysAndValues...)
}

func (l *SlogWrapper) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, keysAndValues...)
}

func (l *SlogWrapper) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, keysAndValues...)
}

func (l *SlogWrapper) ParseLogLevel(level string) (slog.Level, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.Level(0), fmt.Errorf("unknown log level: %s", level)
	}
}

package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(level, format, output string) *Logger {
	// Parse log level
	logLevel := parseLevel(level)
	zerolog.SetGlobalLevel(logLevel)

	// Configure output
	var writer io.Writer
	switch output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		writer = os.Stdout
	}

	// Configure format
	if format == "pretty" || format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
		}
	}

	// Create logger
	l := zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{logger: &l}
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.Debug().Fields(convertFields(fields)).Msg(msg)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.Info().Fields(convertFields(fields)).Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.Warn().Fields(convertFields(fields)).Msg(msg)
}

func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	l.logger.Error().Err(err).Fields(convertFields(fields)).Msg(msg)
}

func (l *Logger) Fatal(msg string, err error, fields ...interface{}) {
	l.logger.Fatal().Err(err).Fields(convertFields(fields)).Msg(msg)
}

func (l *Logger) With(key string, value interface{}) *Logger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &Logger{logger: &newLogger}
}

func (l *Logger) GetZerolog() *zerolog.Logger {
	return l.logger
}

func convertFields(fields []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		m[key] = fields[i+1]
	}
	return m
}

// Global logger instance
var Global *Logger

func InitGlobal(level, format, output string) {
	Global = New(level, format, output)
	log.Logger = *Global.logger
}

func Debug(msg string, fields ...interface{}) {
	if Global != nil {
		Global.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...interface{}) {
	if Global != nil {
		Global.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...interface{}) {
	if Global != nil {
		Global.Warn(msg, fields...)
	}
}

func Error(msg string, err error, fields ...interface{}) {
	if Global != nil {
		Global.Error(msg, err, fields...)
	}
}

func Fatal(msg string, err error, fields ...interface{}) {
	if Global != nil {
		Global.Fatal(msg, err, fields...)
	}
}

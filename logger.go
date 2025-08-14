package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	maxResponseSize = 1024 // 1Kb
)

var logger *zap.Logger

func InitLogger() error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Настройка ротации логов
	logFile := &lumberjack.Logger{
		Filename:   "loadtest.log",
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // дни
		Compress:   true,
	}

	// Форматирование JSON
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(logFile),
			zapcore.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		),
	)

	logger = zap.New(core)
	return nil
}

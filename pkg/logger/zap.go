package logger

import (
	"context"
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/consts"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger

func InitLogger(cfg *config.Config) {
	var (
		core   zapcore.Core
		level  zap.AtomicLevel
		stdout = zapcore.AddSync(os.Stdout)
	)
	if err := level.UnmarshalText([]byte(cfg.Log.Level)); err != nil {
		panic(err)
	}

	if cfg.Env != "" && cfg.Env != "local" {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		file := zapcore.AddSync(&lumberjack.Logger{
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
			Filename:   filepath.Join(cfg.Log.Path, cfg.Log.File),
		})
		fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
		consoleEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
			zapcore.NewCore(fileEncoder, file, level),
		)
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(config.EncoderConfig)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
		)
	}

	zapLog = zap.New(core, zap.WithCaller(true))
}

func Sync() {
	zapLog.Sync()
}

func GetLogger() *zap.Logger {
	return zapLog
}

func WithTrace(ctx context.Context, message string) string {
	traceid := ctx.Value(consts.TraceKey)
	if traceid != nil {
		return fmt.Sprintf("[%s] %s", traceid, message)
	}
	return message
}

func Info(ctx context.Context, message string, fields ...zap.Field) {
	zapLog.WithOptions(zap.AddCallerSkip(1)).Info(WithTrace(ctx, message), fields...) // Add CallerSkip to show true caller instead of logger/zap.go
}

func Debug(ctx context.Context, message string, fields ...zap.Field) {
	zapLog.WithOptions(zap.AddCallerSkip(1)).Debug(WithTrace(ctx, message), fields...)
}

func Error(ctx context.Context, message string, err error) {
	zapLog.WithOptions(zap.AddCallerSkip(1)).Error(WithTrace(ctx, message), zap.Error(err))
}

func Fatal(ctx context.Context, message string, fields ...zap.Field) {
	zapLog.WithOptions(zap.AddCallerSkip(1)).Fatal(WithTrace(ctx, message), fields...)
}

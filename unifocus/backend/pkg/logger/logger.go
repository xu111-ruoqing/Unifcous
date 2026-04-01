package logger

import (
	"os"

	"github.com/unifocus/backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// Init 初始化日志
func Init(cfg *config.LogConfig) error {
	// 日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 输出配置
	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "stdout" {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		// TODO: 实现文件输出和日志轮转
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log = logger.Sugar()

	return nil
}

// Sync 刷新日志缓冲
// 注意: 在应用关闭时调用，确保所有日志都被写入
// 虽然忽略了错误，但在生产环境中可以考虑记录错误或重试
func Sync() {
	if log != nil {
		if err := log.Sync(); err != nil {
			// 在开发环境可以输出错误，生产环境可以忽略或记录到监控系统
			// 因为Sync()通常在应用关闭时调用，此时可能无法再写入日志
		}
	}
}

// Debug 调试日志
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	log.Fatalf(template, args...)
}

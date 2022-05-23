package infra

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kratosx-fashion/app/system/conf"
	"kratosx-fashion/pkg/logutil"
	"os"
	"strings"
	"time"
)

func NewLogger(conf *conf.Logger) log.Logger {
	var zlog *zap.Logger

	if ok, _ := logutil.PathExists(conf.Dir); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", conf.Dir)
		_ = os.Mkdir(conf.Dir, os.ModePerm)
	}
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := []zapcore.Core{
		getEncoderCore(conf, fmt.Sprintf("./%s/debug", conf.Dir), debugPriority),
		getEncoderCore(conf, fmt.Sprintf("./%s/info", conf.Dir), infoPriority),
		getEncoderCore(conf, fmt.Sprintf("./%s/warn", conf.Dir), warnPriority),
		getEncoderCore(conf, fmt.Sprintf("./%s/error", conf.Dir), errorPriority),
	}
	id, _ := os.Hostname()
	switch conf.Level {
	case "info", "INFO":
		zlog = zap.New(
			zapcore.NewTee(cores...),
			zap.Fields(zap.String("service.id", id), zap.String("service.name", conf.Prefix)),
		)
	case "debug", "DEBUG":
		zlog = zap.New(
			zapcore.NewTee(cores...),
			zap.Fields(zap.String("service.id", id), zap.String("service.name", conf.Prefix)),
			zap.AddStacktrace(zap.DebugLevel),
		)
	case "warn", "WARN":
		zlog = zap.New(
			zapcore.NewTee(cores...),
			zap.Fields(zap.String("service.id", id), zap.String("service.name", conf.Prefix)),
			zap.AddStacktrace(zap.WarnLevel),
		)
	case "error", "ERROR":
		zlog = zap.New(
			zapcore.NewTee(cores...),
			zap.Fields(zap.String("service.id", id), zap.String("service.name", conf.Prefix)),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	default:
		zlog = zap.New(
			zapcore.NewTee(cores...),
			zap.Fields(zap.String("service.id", id), zap.String("service.name", conf.Prefix)),
		)
	}

	if conf.ShowLine {
		zlog = zlog.WithOptions(zap.AddCaller())
	}

	return logutil.NewLogger(zlog)
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(conf *conf.Logger) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  conf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder(conf),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case conf.EncodeLevel == "Lowercase": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case conf.EncodeLevel == "LowercaseColor": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case conf.EncodeLevel == "Capital": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case conf.EncodeLevel == "CapitalColor": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder(conf *conf.Logger) zapcore.Encoder {
	if conf.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(conf))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(conf))
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(conf *conf.Logger, fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer, err := logutil.GetWriteSyncer(fileName, conf.LogInConsole) // 使用file-rotatelogs进行日志分割
	if err != nil {
		panic(err)
	}
	return zapcore.NewCore(getEncoder(conf), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(conf *conf.Logger) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(
			t.Format("2006/01/02 - 15:04:05.000" + "\t" + "[" + strings.ToUpper(conf.Prefix) + "]"))
	}
}

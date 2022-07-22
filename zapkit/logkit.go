package zapkit

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logkit *Logkit

type Logkit struct {
	writer *zap.Logger
	sugar  *zap.SugaredLogger
	config *ZapkitConfig
	level  int
}

type ZapkitConfig struct {
	File       string `yaml:"file"`
	Level      string `yaml:"level"`
	MaxSize    int    `yaml:"maxsize"`
	MaxBackups int    `yaml:"maxbackups"`
	MaxAge     int    `yaml:"maxage"`
	Compress   bool   `yaml:"compress"`
}

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone
)

var LoggerLevel = map[string]int{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"none":  LevelNone,
}

// 初始化方法
func Init(zapCfg *ZapkitConfig, args ...string) error {
	if logkit != nil {
		return errors.New("already initialized")
	}
	if zapCfg.File == "" {
		zapCfg.File = "/tmp/zapkit.log"
	}
	if zapCfg.Level == "" {
		zapCfg.Level = "info"
	}
	if zapCfg.MaxSize <= 0 {
		zapCfg.MaxSize = 512
	}
	if zapCfg.MaxBackups <= 0 {
		zapCfg.MaxBackups = 10
	}
	if zapCfg.MaxAge <= 0 {
		zapCfg.MaxAge = 7
	}

	extName := ""
	if len(args) > 0 {
		extName = args[0]
	}

	return initZapkit(zapCfg, extName)
}

func initZapkit(zapCfg *ZapkitConfig, extName string) error {
	logName := zapCfg.File
	if extName != "" {
		logExt := filepath.Ext(zapCfg.File)
		logName = strings.TrimRight(zapCfg.File, logExt) + "_" + extName + logExt
	}

	hook := lumberjack.Logger{
		Filename:   logName,
		MaxSize:    zapCfg.MaxSize,    // megabytes
		MaxBackups: zapCfg.MaxBackups, // 最大保留备份文件个数
		MaxAge:     zapCfg.MaxAge,     // days
		Compress:   zapCfg.Compress,   // 是否压缩，disabled by default
	}

	w := zapcore.AddSync(&hook)

	var level zapcore.Level
	logLevel := strings.ToLower(zapCfg.Level)
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.ErrorLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)

	iLevel, ok := LoggerLevel[logLevel]
	if !ok {
		iLevel = 3
	}
	logger := zap.New(core)
	sugar := logger.Sugar()
	logkit = &Logkit{
		writer: logger,
		sugar:  sugar,
		config: zapCfg,
		level:  iLevel,
	}

	return nil
}

func Debug(str string, args ...zap.Field) {
	logkit.writer.Debug(str, args...)
}

func Info(str string, args ...zap.Field) {
	logkit.writer.Info(str, args...)
}

func Warn(str string, args ...zap.Field) {
	logkit.writer.Warn(str, args...)
}

func Error(str string, args ...zap.Field) {
	logkit.writer.Error(str, args...)
}

func Debugf(str string, args ...interface{}) {
	if logkit.level > LevelDebug {
		return
	}
	logkit.sugar.Debugf(fmt.Sprintf(str, args...))
}
func Infof(str string, args ...interface{}) {
	if logkit.level > LevelInfo {
		return
	}
	logkit.sugar.Infof(fmt.Sprintf(str, args...))
}
func Warnf(str string, args ...interface{}) {
	if logkit.level > LevelWarn {
		return
	}
	logkit.sugar.Warnf(fmt.Sprintf(str, args...))
}
func Errorf(str string, args ...interface{}) {
	if logkit.level > LevelError {
		return
	}
	logkit.sugar.Errorf(str, args...)
}

func Sync() error {
	return logkit.writer.Sync()
}

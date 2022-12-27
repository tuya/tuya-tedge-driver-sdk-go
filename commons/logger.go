package commons

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
)

type TedgeLogger interface {
	SetLogLevel(level string)
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type LoggerConfig struct {
	FileName   string
	Prefix     string // service name
	LogLevel   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

var LogLevelMap = map[string]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

type defaultLogger struct {
	level     *zap.AtomicLevel
	zapLogger *zap.SugaredLogger
}

func DefaultLogger(logLevel, serviceName string) TedgeLogger {
	fileName := "/mnt/logs/driver.log"
	dirBase := filepath.Dir(fileName)
	_, fileErr := os.Stat(dirBase)
	if fileErr != nil {
		os.MkdirAll(dirBase, 0755)
	}

	_, ok := LogLevelMap[logLevel]
	if !ok {
		logLevel = InfoLevel
	}

	return NewDefaultLogger(fileName, logLevel, serviceName)
}

func NewDefaultLogger(fileName, logLevel, serviceName string) TedgeLogger {
	newCfg := &LoggerConfig{
		FileName:   fileName,
		LogLevel:   logLevel,
		Prefix:     fmt.Sprintf("[%s]", serviceName),
		MaxSize:    30, // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 10, // 日志文件最多保存多少个备份
		MaxAge:     7,  // 文件最多保存多少天
		Compress:   false,
	}
	return initLogger(newCfg)
}

func initLogger(lc *LoggerConfig) TedgeLogger {
	var ll zapcore.Level
	if err := ll.UnmarshalText([]byte(lc.LogLevel)); err != nil {
		ll = zapcore.InfoLevel
	}
	var level = zap.NewAtomicLevelAt(ll)
	if lc.FileName == "" {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = level
		cfg.EncoderConfig.ConsoleSeparator = " "
		cfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ := cfg.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.PanicLevel))
		return &defaultLogger{
			level:     &level,
			zapLogger: logger.Sugar().Named(lc.Prefix),
		}
	}

	writeSyncer := getLogWriter(lc)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, level.Level())
	logger := zap.New(core, zap.AddCallerSkip(1), zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))
	return &defaultLogger{
		level:     &level,
		zapLogger: logger.Sugar().Named(lc.Prefix),
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(cfg *LoggerConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func (dl *defaultLogger) SetLogLevel(level string) {
	v, ok := LogLevelMap[level]
	if !ok {
		return
	}
	dl.level.SetLevel(v)
}

func (dl *defaultLogger) Debugf(template string, args ...interface{}) {
	dl.zapLogger.Debugf(template, args...)
}

func (dl *defaultLogger) Infof(template string, args ...interface{}) {
	dl.zapLogger.Infof(template, args...)
}

func (dl *defaultLogger) Warnf(template string, args ...interface{}) {
	dl.zapLogger.Warnf(template, args...)
}

func (dl *defaultLogger) Errorf(template string, args ...interface{}) {
	dl.zapLogger.Errorf(template, args...)
}

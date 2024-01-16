package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Zap = new(_zap)

type _zap struct{}

var Client *zap.Logger

// new logger实列
func NewZap(conf ZapConf) *zap.Logger {
	// 获取cores
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(conf.Level))
	if err != nil {
		panic(err.Error())
	}
	encoder := Zap.getEncoder()
	// 获取日志写入位置
	writeSyncer := Zap.getWriteSyncer(conf.Director, conf.MaxSize, conf.MaxBackup, conf.MaxAge)
	core := zapcore.NewCore(encoder, writeSyncer, l)
	Client = zap.New(core, zap.AddCaller())
	return Client
}

// 日志分割
func (z *_zap) getWriteSyncer(director string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   director,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GetEncoderConfig 获取zapcore.EncoderConfig
func (z *_zap) getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktraceKey",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     z.customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// CustomTimeEncoder 自定义日志输出时间格式
func (z *_zap) customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}

// GetEncoder 获取 zapcore.Encoder
func (z *_zap) getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(z.getEncoderConfig())
}

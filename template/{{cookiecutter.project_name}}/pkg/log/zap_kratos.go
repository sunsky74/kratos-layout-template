package pkg

import (
	"fmt"
	"{{cookiecutter.project_name}}/configs/conf"
	"time"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg *conf.Zap, g *conf.Global) (log.Logger, error) {

	writer := NewLoggerWriter(fmt.Sprintf("%s/%s-%s.log", cfg.Filename, g.AppName, time.Now().Format("20060102")), //文件名
		int(cfg.MaxSize),
		int(cfg.MaxBackups),
		int(cfg.MaxAge),
		cfg.Compress,
		true)
	return newLogger(g, writer)
}

func newLogger(g *conf.Global, w *lumberjack.Logger) (log.Logger, error) {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger := NewZapLogger(encoder, zapcore.AddSync(w), zapcore.DebugLevel,
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development())

	l := kzap.NewLogger(logger)

	return log.With(
		l,
		"env", g.Env,
		"service_id", g.Id,
		"service_name", g.AppName,
		"service_version", g.Version,
		"trace_id", tracing.TraceID(),
	), nil
}

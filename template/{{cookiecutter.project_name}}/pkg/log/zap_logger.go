package pkg

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(encoder zapcore.EncoderConfig, logWrite zapcore.WriteSyncer, level zapcore.Level, opts ...zap.Option) *zap.Logger {

	syncers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	if logWrite != nil {
		syncers = append(syncers, logWrite)
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(syncers...),
		level,
	)
	return zap.New(core, opts...)
}

type LoggerWriterOption func(logger *lumberjack.Logger)

func NewLoggerWriter(filename string, maxSize, maxBackups, maxAge int, compress bool, local bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  local,
	}
}

func NewDefaultLoggerWriter() *lumberjack.Logger {
	return NewLoggerWriterWithOptions(
		WithFilename("./logs/app.log"),
		WithMaxSize(100),
		WithMaxBackups(7),
		WithMaxAge(30),
		WithCompress(true),
		WithLocalTime(true))
}

func NewLoggerWriterWithOptions(opts ...LoggerWriterOption) *lumberjack.Logger {
	if len(opts) == 0 {
		return NewDefaultLoggerWriter()
	}
	w := &lumberjack.Logger{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func WithFilename(name string) LoggerWriterOption {
	return func(w *lumberjack.Logger) {
		w.Filename = name
	}
}

func WithMaxSize(size int) LoggerWriterOption {
	return func(w *lumberjack.Logger) {
		w.MaxSize = size
	}
}

func WithMaxBackups(backups int) LoggerWriterOption {
	return func(w *lumberjack.Logger) {
		w.MaxBackups = backups
	}
}

func WithMaxAge(age int) LoggerWriterOption {
	return func(w *lumberjack.Logger) {
		w.MaxAge = age
	}
}

func WithCompress(compress bool) LoggerWriterOption {
	return func(w *lumberjack.Logger) {
		w.Compress = compress
	}
}

func WithLocalTime(local bool) LoggerWriterOption {
	return func(w *lumberjack.Logger) {
		w.LocalTime = local
	}
}

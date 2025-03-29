package log

import (
	"context"
	"github.com/DavidMovas/gopherbox/pkg/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type Logger struct {
	*zap.Logger
	file   io.Closer
	closer *closer.Closer
}

func NewLogger(level zapcore.Level, local bool) (*Logger, error) {
	logger := &Logger{}

	c := closer.NewCloser()
	logger.closer = c

	var cfg zap.Config
	if local {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.DisableStacktrace = true
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.OutputPaths = []string{"stdout"}
	logger.Logger, _ = cfg.Build(zap.WithCaller(true))

	c.Push(logger.Logger.Sync)

	return logger, nil
}

func (l *Logger) Zap() *zap.Logger {
	return l.Logger
}

func (l *Logger) Stop() error {
	return l.closer.Close(context.Background())
}

package log

import (
	"context"
	"github.com/Brain-Wave-Ecosystem/go-common/pkg/abstractions"
	"io"

	"github.com/DavidMovas/gopherbox/pkg/closer"
	"go.uber.org/zap"
)

var _ abstractions.ILogger = (*Logger)(nil)

type Logger struct {
	zap    *zap.Logger
	file   io.Closer
	closer *closer.Closer
}

func NewLogger(local bool, level string) *Logger {
	logger := &Logger{}

	c := closer.NewCloser()
	logger.closer = c

	atomicLevel := levelFromString(level)

	var cfg zap.Config
	if local {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.DisableStacktrace = true
	cfg.Level = atomicLevel
	cfg.OutputPaths = []string{"stdout"}
	logger.zap, _ = cfg.Build(zap.WithCaller(true))

	c.Push(logger.zap.Sync)

	return logger
}

func (l *Logger) Close() error {
	return l.closer.Close(context.Background())
}

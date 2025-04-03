package log

import (
	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap/zapcore"
)

func (l *Logger) HCLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       "hclog",
		Level:      fromZapLevel(l.zap.Level()),
		JSONFormat: true,
	})
}

func fromZapLevel(lvl zapcore.Level) hclog.Level {
	switch lvl {
	case zapcore.DebugLevel:
		return hclog.Debug
	case zapcore.InfoLevel:
		return hclog.Info
	case zapcore.WarnLevel:
		return hclog.Warn
	case zapcore.ErrorLevel:
		return hclog.Error
	default:
		return hclog.Info
	}
}

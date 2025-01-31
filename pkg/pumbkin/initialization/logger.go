package initialization

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"simpledatabase/pkg/pumbkin/config"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

const (
	defaultEncoding   = "json"
	defaultLevel      = zapcore.InfoLevel
	defaultOutputPath = "spider.log"
)

var supportedLoggingLevels = map[string]zapcore.Level{
	debugLevel: zapcore.DebugLevel,
	infoLevel:  zapcore.InfoLevel,
	warnLevel:  zapcore.WarnLevel,
	errorLevel: zapcore.ErrorLevel,
}

func CreateLogger(cfg *config.Logging) (*zap.Logger, error) {
	level := defaultLevel
	output := defaultOutputPath

	if cfg != nil {
		if cfg.Level != "" {
			var found bool
			if level, found = supportedLoggingLevels[cfg.Level]; !found {
				return nil, errors.New("logging level is incorrect")
			}
		}
		if cfg.Output != "" {
			output = cfg.Output
		}
	}
	loggerCfg := zap.Config{
		Encoding:    defaultEncoding,
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{output},
	}
	return loggerCfg.Build()
}

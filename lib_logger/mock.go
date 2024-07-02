package lib_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func MockLogs() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	return zap.New(core), logs
}

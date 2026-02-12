package logger

import "go.uber.org/zap"

var global *zap.Logger

func SetGlobal(l *zap.Logger) {
	global = l
	zap.ReplaceGlobals(l)
}

func L() *zap.Logger {
	if global == nil {
		panic("global logger not set")
	}
	return global
}

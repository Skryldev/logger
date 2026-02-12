package logger

import "go.uber.org/zap/zapcore"

type Config struct {
	Env          string
	Service      string
	Level        zapcore.Level
	JSON         bool
	Sampling     bool
	EnableCaller bool
}


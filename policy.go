package logger

	import (
		"go.uber.org/zap"
		"go.uber.org/zap/zapcore"
	)

func LevelFromStatus(status int) zapcore.Level {
	switch {
	case status >= 500:
		return zap.ErrorLevel
	case status >= 400:
		return zap.WarnLevel
	default:
		return zap.InfoLevel
	}
}
func main() {
	LevelFromStatus(500)
}
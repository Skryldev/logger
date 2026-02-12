package logger

import (
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

	func DevTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		loc, err := time.LoadLocation("Asia/Tehran")
		if err != nil {
			loc = time.FixedZone("Iran", 3*3600+30*60)
		}
		now := t.In(loc)
		enc.AppendString(now.Format("15:04:05"))
	}

	func ConsoleLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var levelColor string
		switch l {
		case zapcore.DebugLevel:
			levelColor = "\033[36mDEBUG\033[0m"
		case zapcore.InfoLevel:
			levelColor = "\033[32mINFO\033[0m"
		case zapcore.WarnLevel:
			levelColor = "\033[33mWARN\033[0m"
		case zapcore.ErrorLevel:
			levelColor = "\033[31mERROR\033[0m"
		case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
			levelColor = "\033[35mFATAL\033[0m"
		default:
			levelColor = l.String()
		}
		enc.AppendString(levelColor)
	}

func New(cfg Config) (*zap.Logger, error) {
	if cfg.Service == "" {
		return nil, errors.New("logger: Service name must be provided")
	}

	var zapCfg zap.Config

	// 1. Preset
	switch cfg.Env {
	case "production":
		zapCfg = zap.NewProductionConfig()
	case "development":
		zapCfg = zap.NewDevelopmentConfig()
	default:
		return nil, errors.New("logger: invalid Env value")
	}

	// 2. Level (single source of truth)
	zapCfg.Level = zap.NewAtomicLevelAt(cfg.Level)

	// 3. Encoding 
	if cfg.JSON {
		zapCfg.Encoding = "json"
	}

	// 4. Sampling
	if cfg.Sampling && cfg.Env == "production" {
		zapCfg.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}

	// 5. Outputs
	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stderr"}

	// 6. Encoder config by env
	enc := &zapCfg.EncoderConfig

	if cfg.Env == "development" {
		enc.EncodeLevel = ConsoleLevelEncoder
		enc.EncodeTime = DevTimeEncoder
		enc.EncodeDuration = zapcore.StringDurationEncoder
	} else {
		enc.TimeKey = "@timestamp"
		enc.LevelKey = "level"
		enc.NameKey = "logger"
		enc.CallerKey = "caller"
		enc.MessageKey = "message"
		enc.StacktraceKey = "stacktrace"
		enc.FunctionKey = ""

		enc.EncodeTime = zapcore.ISO8601TimeEncoder
		enc.EncodeLevel = zapcore.LowercaseLevelEncoder
		enc.EncodeDuration = zapcore.SecondsDurationEncoder
	}

	enc.EncodeCaller = zapcore.ShortCallerEncoder
	enc.EncodeName = zapcore.FullNameEncoder

	// 7. Options
	opts := []zap.Option{
		zap.Fields(zap.String("service", cfg.Service)),
	}

	if cfg.EnableCaller {
		opts = append(opts, zap.AddCaller())
	}

	return zapCfg.Build(opts...)
}


package logger

import (
	"context"

	"go.uber.org/zap"
)

type key struct{}

func Inject(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, key{}, l)
}

func From(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(key{}).(*zap.Logger); ok {
		return l
	}
	return L()
}

package middleware

import (
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/askari/gpm/logger"
)

type Middleware struct {
	Base *zap.Logger
}

func New(base *zap.Logger) *Middleware {
	return &Middleware{Base: base}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		l := m.Base.With(
			zap.String("request_id", reqID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("client_ip", ip),
		)

		ctx := logger.Inject(r.Context(), l)
		r = r.WithContext(ctx)

		rw := NewResponseWriter(w)

		defer func() {
			if rec := recover(); rec != nil {
				l.Error(
					"panic_recovered",
					zap.Any("error", rec),
					zap.ByteString("stack", debug.Stack()),
				)
				http.Error(rw, "internal error", http.StatusInternalServerError)
				return
			}

			latency := time.Since(start)
			level := logger.LevelFromStatus(rw.status)

			l.Check(level, "http_request_completed").
				Write(
					zap.Int("status", rw.status),
					zap.Duration("latency", latency),
				)
		}()

		next.ServeHTTP(rw, r)
	})
}


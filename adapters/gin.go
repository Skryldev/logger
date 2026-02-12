package adapters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/askari/gpm/logger/middleware"
)

func Logger(l *zap.Logger) gin.HandlerFunc {
	mw := middleware.New(l)

	return func(c *gin.Context) {
		mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

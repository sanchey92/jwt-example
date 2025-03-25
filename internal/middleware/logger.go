package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/sanchey92/jwt-example/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.size += size
	return l.size, err
}

func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetLogger()
			start := time.Now()

			lrw := &loggingResponseWriter{ResponseWriter: w}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			log.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", lrw.statusCode),
				zap.Duration("duration", duration),
				zap.String("remote_ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()))
		})
	}
}

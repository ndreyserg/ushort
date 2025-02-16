// Модуль с функциями для  логгирования
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write - обертка над методом Write для ResponseWriter
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader - обертка над методом WriteHeader для ResponseWriter
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Log - инстанс логгера
var Log zap.SugaredLogger = *zap.NewNop().Sugar()

// Initialize инициализация объекта логгера
func Initialize(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)

	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		return err
	}
	defer zl.Sync()

	Log = *zl.Sugar()

	return nil
}

// LoggerMiddleware middleware для логирования общей информации о запросе
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		start := time.Now()
		method := r.Method
		path := r.URL.Path
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		Log.Infoln(
			"method", method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
			"path", path,
		)
	})
}

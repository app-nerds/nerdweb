package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type requestLogger struct {
	handler http.Handler
	logger  *logrus.Entry
}

func (m *requestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recorder := &statusRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}

	startTime := time.Now()
	ip := realIP(r)

	m.handler.ServeHTTP(recorder, r)
	diff := time.Since(startTime)

	m.logger.WithFields(logrus.Fields{
		"ip":            ip,
		"method":        r.Method,
		"status":        recorder.Status,
		"executionTime": diff,
		"queryParams":   r.URL.RawQuery,
	}).Info(r.URL.Path)
}

/*
RequestLogger returns a middleware for logging all requests.

Example:

  mux := nerdweb.NewServeMux()
  mux.HandleFunc("/endpoint", handler)

  mux.Use(middlewares.RequestLogger(logger))
*/
func RequestLogger(logger *logrus.Entry) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := &requestLogger{
				handler: next,
				logger:  logger,
			}

			handler.ServeHTTP(w, r)
		})
	}
}

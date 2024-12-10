package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logrus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(w, r)

		logrus.WithFields(logrus.Fields{
			"path":     r.URL.Path,
			"method":   r.Method,
			"status":   wrapper.StatusCode,
			"duration": time.Since(start).Milliseconds(),
		}).Info("HTTP request")
	})
}

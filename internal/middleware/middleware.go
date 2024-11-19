package middleware

import (
	"fmt"
	"net/http"
	"time"

	"log/slog"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, req)

		slog.Info(fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, time.Since(start)))
	})
}

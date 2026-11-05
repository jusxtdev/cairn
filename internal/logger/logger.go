package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func InitLogger(logfile string) error {
	// logfile is created at logs/
	// make sure logs/ folder exists
	logFolder := "logs"
	logfilePath := fmt.Sprintf("%s/%s", logFolder, logfile)
	file, err := os.OpenFile(
		logfilePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	logger := slog.New(
		slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	slog.SetDefault(logger)
	return nil
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			slog.Debug("incoming request", 
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(w, r)

			slog.Debug("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
}
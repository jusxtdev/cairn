package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logfile string) (*lumberjack.Logger, error) {
	logFolder := "logs"
	if err := os.MkdirAll(logFolder, 0750); err != nil {
        return nil, err
    }

	logfilePath := fmt.Sprintf("%s/%s", logFolder, logfile)
	file := &lumberjack.Logger{
        Filename:   logfilePath,
        MaxSize:    100, // Maximum size in MB
        MaxBackups: 3,   // Keep 3 old log files
        MaxAge:     7,   // Keep old logs for 7 days
        Compress:   true,
    }

	logger := slog.New(
		slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	slog.SetDefault(logger)
	return file, nil
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
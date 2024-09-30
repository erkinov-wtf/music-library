package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"io"
	"log/slog"
	"music-library/pkg/utils/environment/variable"
	"os"
	"time"
)

var Logger *slog.Logger
var Msg string

// SetupLogger initializes the custom logger and returns the logger instance.
func SetupLogger(env string) *slog.Logger {
	logWriter := &lumberjack.Logger{
		Filename:   getLogName(),
		MaxSize:    10,
		MaxAge:     7,
		MaxBackups: 6,
		Compress:   false,
	}

	multiWriter := io.MultiWriter(logWriter, os.Stdout)

	var handler slog.Handler

	switch env {
	case variable.Debug:
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	case variable.Release:
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		panic("wrong environment")
	}

	Logger = slog.New(handler)
	return Logger
}

func getLogName() string {
	date := time.Now().Local()
	return fmt.Sprintf("logs/music_library_%v.log", date)
}

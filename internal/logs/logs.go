package logs

import (
	"log/slog"
	"os"
)

func InitLogger() (*slog.Logger, func()) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {

		// Fallback
		slog.Warn("cannot open log file, using stderr", "error", err)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		return logger, func() {}
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	return logger, func() {
		err = logFile.Close()
		if err != nil {
			logger.Error("cannot close log file", "error", err)
		}
	}
}

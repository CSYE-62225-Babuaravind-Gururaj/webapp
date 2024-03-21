package logs

import (
	"os"

	"github.com/rs/zerolog"
)

func CreateLogger() zerolog.Logger {

	zerolog.LevelFieldName = "severity"

	logFilePath := "logs/app.log"
	if os.Getenv("RUN_ENV") != "test" {
		logFilePath = "/var/log/myapp/app.log"
	}

	// Open or create the log file.
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	}

	// Create a multi-writer to both the file and stdout.
	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	return logger
}

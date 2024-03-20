package logs

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func CreateLogger() *zap.Logger {

	log.SetOutput(os.Stdout)

	if os.Getenv("RUN_ENV") == "test" {
		// Configure logger for test environment: log to stdout only.
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zap.InfoLevel))
		return zap.New(core)
	}

	logDirPath := "/var/log/myapp"
	logFilePath := logDirPath + "/app.log"

	log.Printf("Creating logger")

	_, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error creating/opening log file: %v\n", err)
		return nil
	}

	// Check if the log file or directory does not exist
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		fmt.Println("Log file is not available")
		log.Printf("Log file is not available")
		// Optionally, create the directory and file here if needed
	} else {
		log.Printf("Log file is available or an error occurred checking the file")
		fmt.Println("Log file is available or an error occurred checking the file")
	}

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/var/log/myapp/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}

package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const timeFormat = "02.01.2006 15:04:05.000"

var logger *log.Logger
var logLevel int8

type logWriter struct{}

func (lw *logWriter) Write(bs []byte) (int, error) {
	return fmt.Print(time.Now().Format(timeFormat), " | ", string(bs))
}

// InitLogger - creates global logger
func InitLogger(serviceName, loggingLevel string) {
	if serviceName == "" {
		serviceName = "UNKNOWN"
	}
	switch loggingLevel {
	case "DEBUG":
		logLevel = 1
	case "ERROR":
		logLevel = 0
	default:
		logLevel = 0
	}

	format := fmt.Sprintf("[%s Service]: ", serviceName)
	logger = log.New(os.Stdout, format, log.Lshortfile)
	logger.SetFlags(0)
	logger.SetOutput(new(logWriter))
}

// Debug - creates debug log message
func Debug(message string) {
	if logLevel >= 1 {
		logger.Printf("DEBUG %s", message)
	}
}

// Error - creates error log message
func Error(message string) {
	logger.Printf("ERROR %s", message)
}

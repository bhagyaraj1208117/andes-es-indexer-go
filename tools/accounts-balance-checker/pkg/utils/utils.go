package utils

import (
	"time"

	logger "github.com/bhagyaraj1208117/andes-logger-go"
)

func LogExecutionTime(log logger.Logger, start time.Time, message string) {
	duration := time.Since(start).Seconds()
	if duration < 0.3 {
		log.Trace(message, "duration in seconds", duration)
	} else {
		log.Debug(message, "duration in seconds", duration)
	}
}

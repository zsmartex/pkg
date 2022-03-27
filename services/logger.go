package services

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zsmartex/pkg/services/logger"
)

func NewLoggerService(service string) *log.Entry {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logger.Formatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	switch os.Getenv("LOG_LEVEL") {
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	}

	return log.WithFields(log.Fields{
		service: service,
	})
}

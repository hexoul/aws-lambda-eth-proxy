// Package log is wrapper for logging
package log

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	// Initalize logger
	logger = log.New()

	// Default configuration
	logger.Formatter = &log.TextFormatter{}
	logger.Out = os.Stdout
	logger.SetLevel(log.WarnLevel)

	// Advanced configuration
	var err error
	var logPath string
	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-log_out" {
			logPath = arg[1]
		} else if arg[0] == "-log_lev" {
			switch strings.ToLower(arg[1]) {
			case "debug":
				logger.SetLevel(log.DebugLevel)
				break
			case "info":
				logger.SetLevel(log.InfoLevel)
				break
			case "warn":
				logger.SetLevel(log.WarnLevel)
				break
			case "error":
				logger.SetLevel(log.ErrorLevel)
				break
			case "fatal":
				logger.SetLevel(log.FatalLevel)
				break
			case "panic":
				logger.SetLevel(log.PanicLevel)
				break
			}
		}
	}
	if logPath != "" {
		if logger.Out, err = os.Create(logPath); err != nil {
			panic("Failed to create log file")
		}
	}
}

// Debug level logging
func Debug(args ...interface{}) {
	logger.Debug(args)
}

// Debugf debug-level logging with format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
}

// Info level logging
func Info(args ...interface{}) {
	logger.Info(args)
}

// Infof info-level logging with format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
}

// Warn level logging
func Warn(args ...interface{}) {
	logger.Warn(args)
}

// Warnf warn-level logging with format
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args)
}

// Error level logging
func Error(args ...interface{}) {
	logger.Error(args)
}

// Errorf error-level logging with format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args)
}

// Fatal level logging and os.Exit
func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

// Fatalf fatal-level logging with format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args)
}

// Panic level logging and panic
func Panic(args ...interface{}) {
	logger.Panic(args)
}

// Panicf panic-level logging with format
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args)
}

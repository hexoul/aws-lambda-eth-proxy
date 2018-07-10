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
		} else if arg[0] == "-log_fmt" {
			switch strings.ToLower(arg[1]) {
			case "text":
				logger.Formatter = &log.TextFormatter{}
				break
			case "json":
				logger.Formatter = &log.JSONFormatter{}
				break
			}
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

// Debugd level logging with id
func Debugd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Debug(args)
}

// Debugf debug-level logging with format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
}

// Debugfd debug-level logging with format and id
func Debugfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Debugf(format, args)
}

// Info level logging
func Info(args ...interface{}) {
	logger.Info(args)
}

// Infod level logging with id
func Infod(id uint64, args ...interface{}) {
	logger.WithField("id", id).Info(args)
}

// Infof info-level logging with format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
}

// Infofd info-level logging with format and id
func Infofd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Infof(format, args)
}

// Warn level logging
func Warn(args ...interface{}) {
	logger.Warn(args)
}

// Warnd level logging with id
func Warnd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Warn(args)
}

// Warnf warn-level logging with format
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args)
}

// Warnfd warn-level logging with format and id
func Warnfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Warnf(format, args)
}

// Error level logging
func Error(args ...interface{}) {
	logger.Error(args)
}

// Errord level logging with id
func Errord(id uint64, args ...interface{}) {
	logger.WithField("id", id).Error(args)
}

// Errorf error-level logging with format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args)
}

// Errorfd error-level logging with format and id
func Errorfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Errorf(format, args)
}

// Fatal level logging and os.Exit
func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

// Fatald level logging with id
func Fatald(id uint64, args ...interface{}) {
	logger.WithField("id", id).Fatal(args)
}

// Fatalf fatal-level logging with format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args)
}

// Fatalfd fatal-level logging with format and id
func Fatalfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Fatalf(format, args)
}

// Panic level logging and panic
func Panic(args ...interface{}) {
	logger.Panic(args)
}

// Panicd level logging with id
func Panicd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Panic(args)
}

// Panicf panic-level logging with format
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args)
}

// Panicfd panic-level logging with format and id
func Panicfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Panicf(format, args)
}

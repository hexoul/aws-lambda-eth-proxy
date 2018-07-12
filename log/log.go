// Package log is wrapper for logging
package log

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var logger, stdoutLogger *log.Logger

func init() {
	// Initalize logger
	logger = log.New()

	// Default configuration
	logger.Formatter = &log.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	}
	logger.Out = os.Stdout
	logger.SetLevel(log.InfoLevel)

	// Advanced configuration
	var err error
	var logPath string
	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-log_out" {
			logPath = arg[1]
			stdoutLogger = log.New()
			stdoutLogger.Out = os.Stdout
		} else if arg[0] == "-log_fmt" {
			switch strings.ToLower(arg[1]) {
			case "json":
				logger.Formatter = &log.JSONFormatter{
					TimestampFormat: "02-01-2006 15:04:05",
				}
				break
			}
		} else if arg[0] == "-log_lev" {
			if lev, err := log.ParseLevel(strings.ToLower(arg[1])); err != nil {
				logger.SetLevel(lev)
			}
		}
	}
	if stdoutLogger != nil {
		stdoutLogger.Formatter = &log.TextFormatter{
			TimestampFormat: "02-01-2006 15:04:05",
			FullTimestamp:   true,
		}
		stdoutLogger.SetLevel(logger.Level)
	}
	if logPath != "" {
		if logger.Out, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
			panic("Failed to create log file")
		}
	}
}

// Debug level logging
func Debug(args ...interface{}) {
	logger.Debug(args)
	if stdoutLogger != nil {
		stdoutLogger.Debug(args)
	}
}

// Debugd level logging with id
func Debugd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Debug(args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Debug(args)
	}
}

// Debugf debug-level logging with format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.Debugf(format, args)
	}
}

// Debugfd debug-level logging with format and id
func Debugfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Debugf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Debugf(format, args)
	}
}

// Info level logging
func Info(args ...interface{}) {
	logger.Info(args)
	if stdoutLogger != nil {
		stdoutLogger.Info(args)
	}
}

// Infod level logging with id
func Infod(id uint64, args ...interface{}) {
	logger.WithField("id", id).Info(args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Info(args)
	}
}

// Infof info-level logging with format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
	if stdoutLogger != nil {
		stdoutLogger.Infof(format, args)
	}
}

// Infofd info-level logging with format and id
func Infofd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Infof(format, args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Infof(format, args)
	}
}

// Warn level logging
func Warn(args ...interface{}) {
	logger.Warn(args)
	if stdoutLogger != nil {
		stdoutLogger.Warn(args)
	}
}

// Warnd level logging with id
func Warnd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Warn(args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Warn(args)
	}
}

// Warnf warn-level logging with format
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.Warnf(format, args)
	}
}

// Warnfd warn-level logging with format and id
func Warnfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Warnf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Warnf(format, args)
	}
}

// Error level logging
func Error(args ...interface{}) {
	logger.Error(args)
	if stdoutLogger != nil {
		stdoutLogger.Error(args)
	}
}

// Errord level logging with id
func Errord(id uint64, args ...interface{}) {
	logger.WithField("id", id).Error(args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Error(args)
	}
}

// Errorf error-level logging with format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.Errorf(format, args)
	}
}

// Errorfd error-level logging with format and id
func Errorfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Errorf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Errorf(format, args)
	}
}

// Fatal level logging and os.Exit
func Fatal(args ...interface{}) {
	logger.Fatal(args)
	if stdoutLogger != nil {
		stdoutLogger.Fatal(args)
	}
}

// Fatald level logging with id
func Fatald(id uint64, args ...interface{}) {
	logger.WithField("id", id).Fatal(args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Fatal(args)
	}
}

// Fatalf fatal-level logging with format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.Fatalf(format, args)
	}
}

// Fatalfd fatal-level logging with format and id
func Fatalfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Fatalf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Fatalf(format, args)
	}
}

// Panic level logging and panic
func Panic(args ...interface{}) {
	logger.Panic(args)
	if stdoutLogger != nil {
		stdoutLogger.Panic(args)
	}
}

// Panicd level logging with id
func Panicd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Panic(args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Panic(args)
	}
}

// Panicf panic-level logging with format
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.Panicf(format, args)
	}
}

// Panicfd panic-level logging with format and id
func Panicfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Panicf(format, args)
	if stdoutLogger != nil {
		stdoutLogger.WithField("id", id).Panicf(format, args)
	}
}

// Package log is wrapper for logging
package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	logger *log.Logger
	// For telegram
	botToken string
	chatID   string
)

func init() {
	// Initalize logger
	logger = log.New()

	// Default configuration
	timestampFormat := "02-01-2006 15:04:05"
	logger.Formatter = &log.TextFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp:   true,
	}
	logger.Out = os.Stdout
	logger.SetLevel(log.InfoLevel)

	// Advanced configuration
	var logPath string
	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-log_out" {
			logPath = arg[1]
		} else if arg[0] == "-log_fmt" {
			switch strings.ToLower(arg[1]) {
			case "json":
				logger.Formatter = &log.JSONFormatter{
					TimestampFormat: timestampFormat,
				}
				break
			}
		} else if arg[0] == "-log_lev" {
			if lev, err := log.ParseLevel(arg[1]); err == nil {
				logger.SetLevel(lev)
			}
		} else if arg[0] == "-log_bot_token" {
			botToken = arg[1]
		} else if arg[0] == "-log_bot_chatid" {
			chatID = arg[1]
		}
	}
	if logPath != "" {
		if f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
			logger.Out = io.MultiWriter(f, os.Stdout)
		} else {
			Panic("Failed to create log file")
		}
		// Stderr
		if f, err := os.OpenFile(logPath+".stderr", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
			redirectStderr(f)
		}
	}
}

func redirectStderr(f *os.File) {
	if err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd())); err != nil {
		Fatalf("Failed to redirect stderr to file: %v", err)
	}
}

func sendTelegramMsg(msg string) {
	if botToken == "" || chatID == "" {
		return
	}
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage?chat_id=" + chatID + "&text=" + msg
	http.Get(url)
}

// Debug level logging
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugd level logging with id
func Debugd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Debug(args...)
}

// Debugf debug-level logging with format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Debugfd debug-level logging with format and id
func Debugfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Debugf(format, args...)
}

// Info level logging
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infod level logging with id
func Infod(id uint64, args ...interface{}) {
	logger.WithField("id", id).Info(args...)
}

// Infof info-level logging with format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Infofd info-level logging with format and id
func Infofd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Infof(format, args...)
}

// Warn level logging
func Warn(args ...interface{}) {
	logger.Warn(args...)
	go sendTelegramMsg(fmt.Sprint("WARN, ", fmt.Sprint(args...)))
}

// Warnd level logging with id
func Warnd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Warn(args...)
	go sendTelegramMsg(fmt.Sprintf("WARN, id:%d, msg:%s", id, fmt.Sprint(args...)))
}

// Warnf warn-level logging with format
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
	go sendTelegramMsg("WARN, " + fmt.Sprintf(format, args...))
}

// Warnfd warn-level logging with format and id
func Warnfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Warnf(format, args...)
	go sendTelegramMsg(fmt.Sprintf("WARN, id:%d, msg:%s", id, fmt.Sprintf(format, args...)))
}

// Error level logging
func Error(args ...interface{}) {
	logger.Error(args...)
	go sendTelegramMsg(fmt.Sprint("ERROR, ", fmt.Sprint(args...)))
}

// Errord level logging with id
func Errord(id uint64, args ...interface{}) {
	logger.WithField("id", id).Error(args...)
	go sendTelegramMsg(fmt.Sprintf("ERROR, id:%d, msg:%s", id, fmt.Sprint(args...)))
}

// Errorf error-level logging with format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
	go sendTelegramMsg("ERROR, " + fmt.Sprintf(format, args...))
}

// Errorfd error-level logging with format and id
func Errorfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Errorf(format, args...)
	go sendTelegramMsg(fmt.Sprintf("ERROR, id:%d, msg:%s", id, fmt.Sprintf(format, args...)))
}

// Fatal level logging and os.Exit
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
	go sendTelegramMsg(fmt.Sprint("FATAL, ", fmt.Sprint(args...)))
}

// Fatald level logging with id
func Fatald(id uint64, args ...interface{}) {
	logger.WithField("id", id).Fatal(args...)
	go sendTelegramMsg(fmt.Sprintf("FATAL, id:%d, msg:%s", id, fmt.Sprint(args...)))
}

// Fatalf fatal-level logging with format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
	go sendTelegramMsg("FATAL, " + fmt.Sprintf(format, args...))
}

// Fatalfd fatal-level logging with format and id
func Fatalfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Fatalf(format, args...)
	go sendTelegramMsg(fmt.Sprintf("FATAL, id:%d, msg:%s", id, fmt.Sprintf(format, args...)))
}

// Panic level logging and panic
func Panic(args ...interface{}) {
	logger.Panic(args...)
	go sendTelegramMsg(fmt.Sprint("PANIC, ", fmt.Sprint(args...)))
}

// Panicd level logging with id
func Panicd(id uint64, args ...interface{}) {
	logger.WithField("id", id).Panic(args...)
	go sendTelegramMsg(fmt.Sprintf("PANIC, id:%d, msg:%s", id, fmt.Sprint(args...)))
}

// Panicf panic-level logging with format
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
	go sendTelegramMsg("PANIC, " + fmt.Sprintf(format, args...))
}

// Panicfd panic-level logging with format and id
func Panicfd(id uint64, format string, args ...interface{}) {
	logger.WithField("id", id).Panicf(format, args...)
	go sendTelegramMsg(fmt.Sprintf("PANIC, id:%d, msg:%s", id, fmt.Sprintf(format, args...)))
}

// Package log is wrapper for logging
package log

import (
	"testing"
)

func TestGeneral(t *testing.T) {
	// Ascending log level
	Debug("debug")
	Info("info")
	Error("error")
	Warn("warn")
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered in ", r)
		}
	}()

	Panic("panic")
}

func TestFatal(t *testing.T) {

	Fatal("fatal")
}

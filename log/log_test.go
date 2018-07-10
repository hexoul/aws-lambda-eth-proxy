// Package log is wrapper for logging
package log

import (
	"testing"

	"github.com/hexoul/aws-lambda-eth-proxy/common"
)

func TestId(t *testing.T) {
	id := common.RandomUint64()
	// Ascending log level
	Debugd(id, "debug", "1")
	Infod(id, "info", "2")
	Warnd(id, "warn", "3")
	Errord(id, "error", "4")
}

func TestGeneral(t *testing.T) {
	// Ascending log level
	Debug("debug", "1")
	Info("info", "2")
	Warn("warn", "3")
	Error("error", "4")
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

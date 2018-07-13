// Package log is wrapper for logging
package log

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hexoul/aws-lambda-eth-proxy/common"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
)

func TestTelegramBot(t *testing.T) {
	accessToken := ""
	chatID := ""
	msg := "alertbot"
	url := "https://api.telegram.org/bot" + accessToken + "/sendMessage?chat_id=" + chatID + "&text=" + msg
	if _, err := http.Get(url); err != nil {
		t.Fatalf("Failed to sendMessage")
	}
}

func TestStderr(t *testing.T) {
	if f, err := os.OpenFile("./test.stderr", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		redirectStderr(f)
	}
	panic("TestStderr")
}

func TestSeq(t *testing.T) {
	f := func(s string) {
		for i := 0; i < 10; i++ {
			Info(s)
		}
	}
	go f("1")
	go f("2")
	go f("3")
	time.Sleep(2 * time.Second)
}

func TestFormatd(t *testing.T) {
	var a uint64
	a = 55555555555555
	Warnf("%d", a)
}

func TestId(t *testing.T) {
	id := common.RandomUint64()
	// Ascending log level
	Debugd(id, "debug", "1")
	Infod(id, "info", "2")
	Warnd(id, "warn", "3")
	Errord(id, "error", "4")
}

func TestJson(t *testing.T) {
	resp := &json.RPCResponse{}
	Error(resp.String())
}

func TestGeneral(t *testing.T) {
	// Ascending log level
	Debug("debug", "1")
	Info("info", "2")
	Warn("warn", "3")
	Error("error", "4")
	time.Sleep(2 * time.Second)
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

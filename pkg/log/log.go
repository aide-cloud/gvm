package log

import (
	"log/slog"
	"sync"
)

var isPrintEnable = true
var once sync.Once

func SetPrintEnable(enable bool) {
	once.Do(func() {
		isPrintEnable = enable
	})
}

func Info(msg string, args ...any) {
	if isPrintEnable {
		slog.Info(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if isPrintEnable {
		slog.Error(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if isPrintEnable {
		slog.Warn(msg, args...)
	}
}

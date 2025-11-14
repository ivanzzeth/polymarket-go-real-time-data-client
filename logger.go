package polymarketrealtime

import "fmt"

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

func (logger) Debug(format string, args ...any) {
	fmt.Printf("[PolymarketRealTimeDataClient][DEBUG] "+format+"\n", args...)
}

func (logger) Info(format string, args ...any) {
	fmt.Printf("[PolymarketRealTimeDataClient][INFO] "+format+"\n", args...)
}

func (logger) Warn(format string, args ...any) {
	fmt.Printf("[PolymarketRealTimeDataClient][WARN] "+format+"\n", args...)
}

func (logger) Error(format string, args ...any) {
	fmt.Printf("[PolymarketRealTimeDataClient][ERROR] "+format+"\n", args...)
}

type silentLogger struct {
}

func (silentLogger) Debug(string, ...any) {}

func (silentLogger) Info(string, ...any) {}

func (silentLogger) Warn(string, ...any) {}

func (silentLogger) Error(string, ...any) {}

func NewSilentLogger() Logger {
	return silentLogger{}
}

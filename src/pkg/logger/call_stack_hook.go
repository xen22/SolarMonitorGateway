package logger

import (
	"runtime/debug"

	"github.com/Sirupsen/logrus"
)

// CallStackHook prints the callstack for fatal or panic type log messages.
type CallStackHook struct {
}

// Levels returns the levels to which this hook applies.
func (hook *CallStackHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel, logrus.PanicLevel}
}

// Fire is used to perform the work of the hook.
func (hook *CallStackHook) Fire(entry *logrus.Entry) error {
	debug.PrintStack()
	return nil
}

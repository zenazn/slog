package slog

import (
	"fmt"
	"time"
)

// Root is the root logger, which logs to stdout by default. Instead of directly
// using the Root logger, we suggest forking and using your own child router
// using Bind().
var Root Logger

func init() {
	// TODO: This is pretty gross. Let's DRY this up
	root := &logger{defaultTarget: stdout, context: context}
	root.lcache = &levelCache{
		iCache: make(map[uintptr]Level),
		logger: root,
		rules:  rules{},
	}
	root.tcache = &targetCache{defaultTarget: stdout}
	Root = root
}

var context = map[string]interface{}{
	"$t": Now,
}

var Now fmt.Stringer = now{}

type now struct{}

func (n now) String() string {
	return time.Now().Format(time.RFC3339Nano)
}

var DefaultLevel = LInfo

func (l Level) String() string {
	switch l {
	case LDebug:
		return "DEBUG"
	case LInfo:
		return "INFO"
	case LWarn:
		return "WARN"
	case LError:
		return "ERROR"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// Bind returns a new "child" Logger that additionally binds the given context
// variables.
func Bind(context map[string]interface{}) Logger {
	return Root.Bind(context)
}

// Set the log level for a given selector on the root logger. See the
// documentation for Logger.SetLevel for the syntax accepted for the selector.
func SetLevel(selector string, level Level) {
	Root.SetLevel(selector, level)
}

// LogTo logs pre-formatted log lines at the given levels to a channel. If you
// do not pass any levels, the channel will be used as the default logger for
// levels not otherwise configured.
func LogTo(target chan<- string, levels ...Level) {
	Root.LogTo(target, levels...)
}

// SlogTo logs raw log lines at the given levels to a channel. If you do not
// pass any levels, the channel will be used as the default logger for levels
// not otherwise configured.
func SlogTo(target chan<- map[string]interface{}, levels ...Level) {
	Root.SlogTo(target, levels...)
}

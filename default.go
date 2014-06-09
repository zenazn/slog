package slog

import (
	"fmt"
	"time"
)

var root logger

func makeRoot(now fmt.Stringer) logger {
	root := logger{
		defaultTarget: stdout,
		context: map[string]interface{}{
			"$time": now,
		},
	}
	root.genLCache(nil)
	root.genTCache(nil)
	return root
}

type currentTime struct{}

func (_ currentTime) String() string {
	// TODO(carl): Just hardcode this. It'll be faster and we can even make
	// it fixed-width.
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func init() {
	root = makeRoot(currentTime{})
}

var DefaultLevel = LInfo

// String implements the fmt.Stringer interface by returning one of DEBUG, INFO,
// WARN, or ERROR.
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
		// Unclear how this would happen, but it's probably not nice to
		// panic().
		return fmt.Sprintf("Level(%d)", l)
	}
}

// Bind returns a new Logger forked from the global root logger that
// additionally binds the given context variables.
func Bind(context map[string]interface{}) Logger {
	return root.Bind(context)
}

// SetLevel sets the log level for a given selector on the root logger. See the
// documentation for Logger.SetLevel for the syntax accepted for the selector.
func SetLevel(selector string, level Level) {
	root.SetLevel(selector, level)
}

// LogTo logs pre-formatted log lines at the given levels to a channel. If you
// do not pass any levels, the channel will be used as the default logger for
// levels not otherwise configured.
func LogTo(target chan<- string, levels ...Level) {
	root.LogTo(target, levels...)
}

// SlogTo logs raw log lines at the given levels to a channel. If you do not
// pass any levels, the channel will be used as the default logger for levels
// not otherwise configured.
func SlogTo(target chan<- map[string]interface{}, levels ...Level) {
	root.SlogTo(target, levels...)
}

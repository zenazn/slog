/*
Package slog is a hierarchical structured logger that supports runtime
modification of log levels on a per-package and per-function basis.
*/
package slog

type Level int

const (
	LDebug Level = 10 + 10*iota
	LInfo
	LWarn
	LError
)

// Type Data is a convenience type for use in constructing objects for Loggers.
type Data map[string]interface{}

/*
Type Logger is the interface provided by slog's structured logger. It includes
facilities for logging at several levels, the ability to create "child" loggers
that inherit the settings of their parent, and the ability to dynamically alter
log verbosity on a per-package and per-function basis at runtime.
*/
type Logger interface {
	// Log at the debug level. Returns true if the current function is
	// configured to log at the debug level.
	Debug(lines ...map[string]interface{}) bool
	// Log at the standard (info) level. Returns true if the current
	// function is configured to log at the standard level.
	Log(lines ...map[string]interface{}) bool
	// Log at the warn level. Returns true if the current function is
	// configured to log at the warn level.
	Warn(lines ...map[string]interface{}) bool
	// Log at the error level. Returns true if the current function is
	// configured to log at the error level.
	Error(lines ...map[string]interface{}) bool

	// TODO(carl): Should we add Fatal (log and os.Exit(1)) or Panic (log
	// and panic a wrapped Error object or something)?

	// Return a new "child" Logger that, in addition to binding the
	// variables in this logger's context, also binds some number of other
	// variables, with newly-bound variables overwriting old variables with
	// the same name. Binding a nil context returns a new Logger that does
	// not bind any additional variables.
	Bind(context map[string]interface{}) Logger

	// Set the log level of a given selector. Selectors are essentially
	// import paths: "github.com/zenazn/slog" selects this package and all
	// subpackages, for instance. If you only wish to select the specific
	// package referenced, you can end the package name with a trailing
	// period ("github.com/zenazn/slog."), which will cause the selector to
	// match functions in the slog package but not any subpackages of slog.
	// If you wish to select only subpackages, you can leave a trailing
	// slash. You can also specify specific function names:
	// "github.com/zenazn/slog.SetLevel" would affect only logs printed
	// within the SetLevel function, for instance.
	//
	// The set of selectors that will be attempted is the union of this
	// Logger's selectors and the selectors of all parent loggers, with
	// selectors specified in children overriding identically named
	// selectors in parents. This inheritence is dynamic: changing the
	// selector set of a parent logger will cause the logging behavior of
	// all children to change.
	//
	// The level a function will log at is the level of the longest selector
	// that matches.
	SetLevel(selector string, level Level)

	// Write log lines for the given levels to the given channel. Logs
	// written to the channel will be single-line strings without a trailing
	// newline that are formatted in a manner that's suitable for immediate
	// printing. As a special case, if no levels are passed, the channel
	// will be used as a default for levels not otherwise specified.
	LogTo(chan<- string, ...Level)

	// Write unformatted log lines for the given level to the given channel.
	// As a special case, if no levels are passed, the channel will be used
	// as a default for levels not otherwise specified.
	SlogTo(chan<- map[string]interface{}, ...Level)
}

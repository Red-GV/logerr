package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ViaQ/logerr/kverrors"
	"github.com/go-logr/logr"
)

// ErrUnknownLoggerType is returned when trying to perform a *Logger only function
// that is incompatible with logr.Logger interface
var ErrUnknownLoggerType = kverrors.New("unknown error type")

var (
	defaultOutput  io.Writer = os.Stdout
	defautLogLevel           = 0

	mtx    sync.RWMutex
	logger logr.Logger = logr.New(NewLogSink("", os.Stdout, 0, JSONEncoder{}))
)

// Init initializes the logger. This is required to use logging correctly
// component is the name of the component being used to log messages.
// Typically this is your application name.
//
// keyValuePairs are key/value pairs to be used with all logs in the future
func Init(component string, keyValuePairs ...interface{}) {
	InitWithOptions(component, nil, keyValuePairs...)
}

// MustInit calls Init and panics if it returns an error
func MustInit(component string, keyValuePairs ...interface{}) {
	Init(component, keyValuePairs...)
}

// InitWithOptions inits the logger with the provided opts
func InitWithOptions(component string, opts []Option, keyValuePairs ...interface{}) {
	mtx.Lock()
	defer mtx.Unlock()

	ls := NewLogSink(component, defaultOutput, 0, JSONEncoder{}, keyValuePairs...)

	for _, opt := range opts {
		opt(ls)
	}

	// don't lock because we already have a lock
	ll := logr.New(ls)
	useLogger(ll)
}

// MustInitWithOptions calls InitWithOptions and panics if an error is returned
func MustInitWithOptions(component string, opts []Option, keyValuePairs ...interface{}) {
	InitWithOptions(component, opts, keyValuePairs...)
}

// GetLogger returns the root logger used for logging
func GetLogger() logr.Logger {
	return logger
}

// UseLogger bypasses the requirement for Init and sets the logger to l
func UseLogger(l logr.Logger) {
	mtx.Lock()
	defer mtx.Unlock()
	useLogger(l)
}

// useLogger sets the logger to l without mtx.Lock()
// To use mtx.Lock see UseLogger
func useLogger(l logr.Logger) {
	logger = l
}

// V returns an Logger value for a specific verbosity level, relative to
// this Logger.  In other words, V values are additive.  V higher verbosity
// level means a log message is less important.
// V(level uint8) Logger
func V(level int) logr.Logger {
	mtx.RLock()
	defer mtx.RUnlock()
	return logger.V(level)
}

// Info logs a non-error message with the given key/value pairs as context.
//
// The msg argument should be used to add some constant description to
// the log line.  The key/value pairs can then be used to add additional
// variable information.  The key/value pairs should alternate string
// keys and arbitrary values.
func Info(msg string, keysAndValues ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Info(msg, keysAndValues...)
}

// Error logs an error, with the given message and key/value pairs as context.
// It functions similarly to calling Info with the "error" named value, but may
// have unique behavior, and should be preferred for logging errors (see the
// package documentations for more information).
//
// The msg field should be used to add context to any underlying error,
// while the err field should be used to attach the actual error that
// triggered this log line, if present.
func Error(err error, msg string, keysAndValues ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Error(err, msg, keysAndValues...)
}

// WithValues adds some key-value pairs of context to a logger.
// See Info for documentation on how key/value pairs work.
func WithValues(keysAndValues ...interface{}) logr.Logger {
	mtx.RLock()
	defer mtx.RUnlock()
	return logger.WithValues(keysAndValues...)
}

// WithName adds a new element to the logger's name.
// Successive calls with WithName continue to append
// suffixes to the logger's name.  It's strongly recommended
// that name segments contain only letters, digits, and hyphens
// (see the package documentation for more information).
func WithName(name string) logr.Logger {
	mtx.RLock()
	defer mtx.RUnlock()
	return logger.WithName(name)
}

// SetLogLevel sets the output verbosity
func SetLogLevel(v int) error {
	mtx.Lock()
	defer mtx.Unlock()
	switch ls := logger.GetSink().(type) {
	case *LogSink:
		ls.SetVerbosity(v)
	default:
		return kverrors.Add(ErrUnknownLoggerType,
			"logger_type", fmt.Sprintf("%T", logger),
			"expected_type", fmt.Sprintf("%T", &LogSink{}),
		)
	}
	return nil
}

// SetOutput sets the logger output to w if the root logger is *log.Logger
// otherwise it returns ErrUnknownLoggerType
func SetOutput(w io.Writer) error {
	mtx.RLock()
	defer mtx.RUnlock()
	switch ls := logger.GetSink().(type) {
	case *LogSink:
		ls.SetOutput(w)
	default:
		return kverrors.Add(ErrUnknownLoggerType,
			"logger_type", fmt.Sprintf("%T", logger),
			"expected_type", fmt.Sprintf("%T", &LogSink{}),
		)
	}
	return nil
}

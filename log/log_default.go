package log

import "io"

// l defines the default log variable for the global log functions.
var l Logger

//==============================================================================

// Init initializes the default logger to allow usage of the global log
// functions.
func Init(w io.Writer, level func() int, flags int) {
	l.mu.Lock()
	{
		dl := New(w, level, flags)

		l.Logger = dl.Logger
		l.level = dl.level
	}
	l.mu.Unlock()
}

//==============================================================================

// Dev logs trace information for developers.
func Dev(traceID string, funcName string, format string, a ...interface{}) {
	l.DevOffset(traceID, 1, funcName, format, a...)
}

// User logs trace information for users.
func User(traceID string, funcName string, format string, a ...interface{}) {
	l.UserOffset(traceID, 1, funcName, format, a...)
}

// Error logs trace information that are errors.
func Error(traceID string, funcName string, err error, format string, a ...interface{}) {
	l.ErrorOffset(traceID, 1, funcName, err, format, a...)
}

// Fatal logs trace information for users and terminates the app.
func Fatal(traceID string, funcName string, format string, a ...interface{}) {
	l.FatalOffset(traceID, 1, funcName, format, a...)
}

//==============================================================================

// DevOffset logs trace information for developers with a offset option to
// expand the caller level.
func DevOffset(traceID string, offset int, funcName string, format string, a ...interface{}) {
	l.DevOffset(traceID, offset+1, funcName, format, a...)
}

// UserOffset logs trace information for users with a offset option to expand the
// caller level.
func UserOffset(traceID string, offset int, funcName string, format string, a ...interface{}) {
	l.UserOffset(traceID, offset+1, funcName, format, a...)
}

// ErrorOffset logs trace information that are errors with a offset option to
// expand the caller level.
func ErrorOffset(traceID string, offset int, funcName string, err error, format string, a ...interface{}) {
	l.ErrorOffset(traceID, offset+1, funcName, err, format, a...)
}

// FatalOffset logs trace information for users and terminates the app with a
// offset expand the caller level.
func FatalOffset(traceID string, offset int, funcName string, format string, a ...interface{}) {
	l.FatalOffset(traceID, offset+1, funcName, format, a...)
}

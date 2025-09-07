package internal

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
)

type ErrorCode uint

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeNotFound
	ErrorCodeInvalidArgument
	ErrorCodeDatabase
	ErrorCodeConfig
	ErrorCodeServer
)

type Error struct {
	orig  error
	msg   string
	code  ErrorCode
	stack string
}

func New(code ErrorCode, format string, a ...interface{}) *Error {
	return &Error{
		code:  code,
		msg:   fmt.Sprintf(format, a...),
		stack: captureStack(),
	}
}
func Wrap(orig error, code ErrorCode, format string, a ...interface{}) *Error {
	return &Error{
		orig:  orig,
		code:  code,
		msg:   fmt.Sprintf(format, a...),
		stack: captureStack(),
	}
}

func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.orig
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func (e *Error) Stack() string {
	return e.stack
}

// Fields returns zap fields for structured logging.
func (e *Error) Fields() []zap.Field {
	return []zap.Field{
		zap.String("msg", e.msg),
		zap.Error(e.orig),
		zap.Uint("code", uint(e.code)),
		zap.String("stack", e.stack),
	}
}

func captureStack() string {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(3, pcs) // skip runtime + this function + wrapper
	frames := runtime.CallersFrames(pcs[:n])

	var stack string
	for {
		frame, more := frames.Next()
		stack += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return stack
}

func NewConfigError(format string, a ...interface{}) *Error {
	return New(ErrorCodeConfig, format, a...)
}

func WrapDatabaseError(orig error, format string, a ...interface{}) *Error {
	return Wrap(orig, ErrorCodeDatabase, format, a...)
}

func NewInvalidArgError(format string, a ...interface{}) *Error {
	return New(ErrorCodeInvalidArgument, format, a...)
}

func NewNotFoundError(format string, a ...interface{}) *Error {
	return New(ErrorCodeNotFound, format, a...)
}

func WrapServerError(orig error, format string, a ...interface{}) *Error {
	return Wrap(orig, ErrorCodeServer, format, a...)
}

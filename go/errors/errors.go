package errors

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"

	"github.com/nownabe/golink/go/clog"
)

func New(msg string) error {
	return newError(nil, msg)
}

func NewWithoutStack(msg string) error {
	return &wrapped{
		err:   nil,
		msg:   msg,
		stack: nil,
	}
}

func Errorf(format string, args ...any) error {
	return newError(nil, fmt.Sprintf(format, args...))
}

func Wrap(err error, msg string) error {
	return wrap(err, msg)
}

func Wrapf(err error, format string, args ...any) error {
	return wrap(err, fmt.Sprintf(format, args...))
}

func newError(err error, msg string) error {
	return &wrapped{
		err:   err,
		msg:   msg,
		stack: callers(),
	}
}

func wrap(err error, msg string) error {
	if w, ok := err.(*wrapped); ok {
		if w.hasStack() {
			return &wrapped{
				err:   err,
				msg:   msg,
				stack: nil,
			}
		}
	}

	return &wrapped{
		err:   err,
		msg:   msg,
		stack: callers(),
	}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

type wrapped struct {
	err     error
	msg     string
	stack   []uintptr
	context *clog.ErrorContext
}

func (w *wrapped) Error() string {
	if w.err == nil {
		return w.msg
	}
	return w.msg + ": " + w.err.Error()
}

func (w *wrapped) Unwrap() error {
	return w.err
}

func (w *wrapped) ErrorContext() *clog.ErrorContext {
	return w.context
}

func (w *wrapped) Stack() []byte {
	buf := bytes.Buffer{}
	buf.WriteString(w.Error())
	buf.WriteString("\n")
	buf.Write(w.frames())
	return buf.Bytes()
}

func (w *wrapped) frames() []byte {
	if w.stack == nil {
		if ww, ok := w.err.(*wrapped); ok {
			return ww.frames()
		}

		return []byte{}
	}

	buf := bytes.Buffer{}

	frames := runtime.CallersFrames(w.stack)
	for {
		f, ok := frames.Next()
		if !ok {
			break
		}
		buf.WriteString(fmt.Sprintf("%s(...)\n", f.Function))
		buf.WriteString(fmt.Sprintf("\t%s:%d\n", f.File, f.Line))
	}

	return buf.Bytes()
}

func (w *wrapped) hasStack() bool {
	if ww, ok := w.err.(*wrapped); ok {
		return w.stack != nil || ww.hasStack()
	}
	return w.stack != nil
}

func callers() []uintptr {
	const depth = 40
	var pcs [depth]uintptr
	// skip [runtime.Callers, this function, errors' unexported constructors, errors' exported custructors]
	n := runtime.Callers(4, pcs[:])
	return pcs[:n]
}

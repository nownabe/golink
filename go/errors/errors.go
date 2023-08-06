package errors

import (
	"runtime/debug"

	"github.com/nownabe/golink/go/clog"
)

func New(msg string) error {
	return &wrapped{
		err:   nil,
		msg:   msg,
		stack: debug.Stack(),
	}
}

func Wrap(err error, msg string) error {
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
		stack: debug.Stack(),
	}
}

type wrapped struct {
	err     error
	msg     string
	stack   []byte
	context *clog.ErrorContext
}

func (w *wrapped) Error() string {
	if w.err == nil {
		return w.msg
	}
	return w.msg + w.err.Error()
}

func (w *wrapped) ErrorContext() *clog.ErrorContext {
	return w.context
}

func (w *wrapped) Stack() []byte {
	if w.stack == nil {
		return w.stack
	}

	if ww, ok := w.err.(*wrapped); ok {
		return ww.Stack()
	}

	return []byte{}
}

func (w *wrapped) hasStack() bool {
	if ww, ok := w.err.(*wrapped); ok {
		return w.stack != nil || ww.hasStack()
	}
	return w.stack != nil
}

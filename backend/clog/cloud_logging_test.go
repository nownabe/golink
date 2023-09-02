package clog_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func TestSourceHandler_Info(t *testing.T) {
	t.Parallel()

	l, w := newLogger()
	ctx := context.Background()

	pc, _, _, _ := runtime.Caller(0) // This file must be just abeve l.Info(ctx, "text")
	l.Info(ctx, "test")

	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()

	want := expectation{
		"time":     anyStringVal{},
		"severity": "INFO",
		"message":  "test",
		"logging.googleapis.com/sourceLocation": map[string]any{
			"file":     frame.File,
			"line":     float64(frame.Line + 1), // json.Unmarshal use float64 for any type
			"function": frame.Function,
		},
	}

	w.expect(t, want)
}

func TestSourceHandler_Err(t *testing.T) {
	t.Parallel()

	l, w := newLogger()
	ctx := context.Background()

	pc, _, _, _ := runtime.Caller(0)
	l.Err(ctx, errors.New("test"))

	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()

	want := expectation{
		"time":     anyStringVal{},
		"severity": "ERROR",
		"message":  "test",
		"logging.googleapis.com/sourceLocation": map[string]any{
			"file":     frame.File,
			"line":     float64(frame.Line + 1), // json.Unmarshal use float64 for any type
			"function": frame.Function,
		},
	}

	fmt.Println(w.String())

	w.expect(t, want)
}

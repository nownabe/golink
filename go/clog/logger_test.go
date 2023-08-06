package clog_test

import (
	"context"
	"testing"
)

// TODO: test other methods

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	l, w := newLogger()

	ctx := context.Background()
	l.Info(ctx, "test")

	want := expectation{
		"time":                                  anyStringVal{},
		"logging.googleapis.com/sourceLocation": nonNilVal{},
		"severity":                              "INFO",
		"message":                               "test",
	}

	w.expect(t, want)
}

func TestLogger_Warning(t *testing.T) {
	t.Parallel()

	l, w := newLogger()

	ctx := context.Background()
	l.Warning(ctx, "test")

	want := expectation{
		"time":                                  anyStringVal{},
		"logging.googleapis.com/sourceLocation": nonNilVal{},
		"severity":                              "WARNING",
		"message":                               "test",
	}

	w.expect(t, want)
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	l, w := newLogger()

	ctx := context.Background()
	l.Error(ctx, "test")

	want := expectation{
		"time":                                  anyStringVal{},
		"logging.googleapis.com/sourceLocation": nonNilVal{},
		"severity":                              "ERROR",
		"message":                               "test",
	}

	w.expect(t, want)
}

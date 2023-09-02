package clog_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nownabe/golink/backend/clog"
)

type anyVal struct{}
type anyStringVal struct{}
type nonNilVal struct{}

type expectation map[string]any

type stubWriter struct {
	*bytes.Buffer
}

func (w *stubWriter) expect(t *testing.T, e expectation) {
	t.Helper()

	l, err := w.ReadBytes('\n')
	if err != nil {
		panic(err)
	}

	got := map[string]any{}
	if err := json.Unmarshal([]byte(l), &got); err != nil {
		panic(err)
	}

	if len(got) != len(e) {
		t.Errorf("got %d keys, want %d", len(got), len(e))
	}

	for k, wantVal := range e {
		gotVal, ok := got[k]
		if !ok {
			t.Errorf("got no key %q", k)
		}

		switch wantVal.(type) {
		case anyVal:
			continue
		case anyStringVal:
			if _, ok := gotVal.(string); !ok {
				t.Errorf("got[%q] is %T (%#v), want string", k, gotVal, gotVal)
			}
		case nonNilVal:
			if gotVal == nil {
				t.Errorf("got[%q] is %v, want nil", k, gotVal)
			}
		default:
			if !cmp.Equal(gotVal, wantVal) {
				t.Errorf("got[%q] is %#v, want %#v: (-want +got)\n%s\n", k, gotVal, wantVal, cmp.Diff(wantVal, gotVal))
			}
		}
	}
}

func newLogger() (*clog.Logger, *stubWriter) {
	w := &stubWriter{&bytes.Buffer{}}
	return clog.New(w, clog.LevelInfo), w
}

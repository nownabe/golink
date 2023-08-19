package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nownabe/golink/backend/middleware"
)

func TestRecoverer(t *testing.T) {
	t.Parallel()

	h := middleware.NewRecoverer()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("panic") }))
	s := httptest.NewServer(h)
	defer s.Close()

	resp := testRequest(t, s, http.MethodGet, "/")

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status code should be %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestRecoverer_ErrAbortHandler(t *testing.T) {
	t.Parallel()

	defer func() {
		rcv := recover()
		if rcv != http.ErrAbortHandler {
			t.Fatal("http.ErrAbortHandler should not be recovered")
		}
	}()

	h := middleware.NewRecoverer()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, req)
}

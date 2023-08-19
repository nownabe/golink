package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, s *httptest.Server, method, path string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), method, s.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

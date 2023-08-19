package backend_test

/*
func TestHandler(t *testing.T) {
	tests := map[string]struct {
		method  string
		reqURL  string
		status  int
		wantURL string
	}{
		"with name":       {"GET", "https://host/linkname", http.StatusMovedPermanently, "https://example.com/redirected"},
		"without name":    {"GET", "https://host/", http.StatusBadRequest, ""},
		"with empty name": {"GET", "https://host//", http.StatusBadRequest, ""},
	}

	h := backend.NewHandler(&fakeRepo{})

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.reqURL, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			got := w.Result()

			if got.StatusCode != tt.status {
				t.Errorf("StatusCode: got %d, want %d", got.StatusCode, tt.status)
			}

			if got.Header.Get("Location") != tt.wantURL {
				t.Errorf("Location: got %s, want %s", got.Header.Get("Location"), tt.wantURL)
			}
		})
	}
}

type fakeRepo struct{}

func (r *fakeRepo) GetURLAndUpdateStats(ctx context.Context, name string) (*url.URL, error) {
	return url.Parse("https://example.com/redirected")
}
*/

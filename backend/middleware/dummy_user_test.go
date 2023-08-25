package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nownabe/golink/backend/middleware"
)

func userRecordHandler(email, userID *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*email = r.Header.Get("X-Appengine-User-Email")
		*userID = r.Header.Get("X-Appengine-User-Id")
	})
}

func TestDummyUser(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		email      string
		userID     string
		wantEmail  string
		wantUserID string
	}{
		"no user": {
			email:      "",
			userID:     "",
			wantEmail:  "",
			wantUserID: "",
		},
		"with user": {
			email:      "user@example.com",
			userID:     "user-id",
			wantEmail:  "accounts.google.com:user@example.com",
			wantUserID: "accounts.google.com:user-id",
		},
		"with email": {
			email:      "user@example.com",
			userID:     "",
			wantEmail:  "",
			wantUserID: "",
		},
		"with user id": {
			email:      "",
			userID:     "user-id",
			wantEmail:  "",
			wantUserID: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var gotEmail, gotUserID string

			h := middleware.NewDummyUser(tt.email, tt.userID)(userRecordHandler(&gotEmail, &gotUserID))
			s := httptest.NewServer(h)
			defer s.Close()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := http.DefaultClient.Do(req); err != nil {
				t.Fatal(err)
			}

			if gotEmail != tt.wantEmail {
				t.Errorf("Email should be %q, but got %q", tt.wantEmail, gotEmail)
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("User ID should be %q, but got %q", tt.wantUserID, gotUserID)
			}
		})
	}
}

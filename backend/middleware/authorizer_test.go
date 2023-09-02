package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nownabe/golink/backend/golinkcontext"
	"github.com/nownabe/golink/backend/middleware"
)

func userContextRecordHandler(email, userID *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		*email, ok = golinkcontext.UserEmailFrom(r.Context())
		if !ok {
			panic("email not found")
		}

		*userID, ok = golinkcontext.UserIDFrom(r.Context())
		if !ok {
			panic("userID not found")
		}

		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthorizer(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		headers    map[string]string
		wantStatus int
		wantEmail  string
		wantUserID string
	}{
		"email and userID": {
			headers: map[string]string{
				"X-Appengine-User-Email": "accounts.google.com:user@example.com",
				"X-Appengine-User-Id":    "accounts.google.com:user-id",
			},
			wantStatus: http.StatusOK,
			wantEmail:  "user@example.com",
			wantUserID: "user-id",
		},
		"email only": {
			headers: map[string]string{
				"X-Appengine-User-Email": "accounts.google.com:user@example.com",
			},
			wantStatus: http.StatusUnauthorized,
			wantEmail:  "",
			wantUserID: "",
		},
		"userID only": {
			headers: map[string]string{
				"X-Appengine-User-Id": "accounts.google.com:user-id",
			},
			wantStatus: http.StatusUnauthorized,
			wantEmail:  "",
			wantUserID: "",
		},
		"no user": {
			headers:    map[string]string{},
			wantStatus: http.StatusUnauthorized,
			wantEmail:  "",
			wantUserID: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var gotEmail, gotUserID string

			h := middleware.NewAuthorizer()(userContextRecordHandler(&gotEmail, &gotUserID))
			s := httptest.NewServer(h)
			defer s.Close()

			req, err := http.NewRequest(http.MethodGet, s.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status: got %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			if resp.StatusCode != http.StatusOK {
				return
			}

			if gotEmail != tt.wantEmail {
				t.Errorf("email: got %v, want %v", gotEmail, tt.wantEmail)
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("userID: got %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

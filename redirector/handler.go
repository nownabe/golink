package redirector

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware is a middleware.
type Middleware func(http.Handler) http.Handler

// NewHandler returns a new Handler.
func NewHandler(repo Repository, middlewares ...Middleware) http.Handler {
	h := &Handler{}
	h.handler = &redirectHandler{repo: repo}

	for _, m := range middlewares {
		h.handler = m(h.handler)
	}

	return h
}

type Handler struct {
	handler http.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

type redirectHandler struct {
	repo Repository
}

func (h *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")

	if len(path) < 1 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if path[0] == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.repo.GetURLAndUpdateStats(ctx, path[0])
	// TODO: If no existence, redirect to /c/name
	if err != nil {
		// TODO: log: error: failed to get url
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(path) > 1 {
		url.Path = fmt.Sprintf("%s/%s", url.Path, strings.Join(path[1:], "/"))
	}

	http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
}

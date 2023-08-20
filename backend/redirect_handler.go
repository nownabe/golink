package backend

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
)

type redirectHandler struct {
	consolePrefix string
	repo          *repository
}

func newRedirectHandler(repo *repository, consolePrefix string) *redirectHandler {
	return &redirectHandler{
		consolePrefix: consolePrefix,
		repo:          repo,
	}
}

func (h *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")

	if path[0] == "" {
		http.Redirect(w, r, h.consolePrefix, http.StatusMovedPermanently)
		return
	}

	url, err := h.repo.GetURLAndUpdateStats(ctx, path[0])
	if err != nil {
		if errors.Is(err, errDocumentNotFound) {
			http.Redirect(w, r, fmt.Sprintf("%s%s", h.consolePrefix, path[0]), http.StatusTemporaryRedirect)
			return
		}

		err := errors.Wrapf(err, "failed to get url: %s", path[0])
		clog.Err(ctx, err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(path) > 1 {
		url.Path = fmt.Sprintf("%s/%s", url.Path, strings.Join(path[1:], "/"))
	}

	http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
}

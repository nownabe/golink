package backend

import (
	"fmt"
	"net/http"
	"net/url"
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

	golink, err := h.repo.Get(ctx, path[0])
	if err != nil {
		if errors.Is(err, errDocumentNotFound) {
			http.Redirect(w, r, fmt.Sprintf("%s%s", h.consolePrefix, path[0]), http.StatusTemporaryRedirect)
			return
		}

		err := errors.Wrapf(err, "failed to get url for %q", path[0])
		clog.Err(ctx, err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	u, err := url.Parse(golink.URL)
	if err != nil {
		err := errors.Wrapf(err, "failed to parse url (id=%q): %q", path[0], golink.URL)
		clog.Err(ctx, err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(path) > 1 {
		u.Path = fmt.Sprintf("%s/%s", u.Path, strings.Join(path[1:], "/"))
	}

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)

	go h.repo.incrementCount(ctx, golink.Name)
}

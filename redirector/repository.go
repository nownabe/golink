package redirector

import (
	"context"
	"net/url"

	"cloud.google.com/go/firestore"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collectionName = "golinks"

var errDocumentNotFound = errors.NewWithoutStack("document not found")

type Repository interface {
	GetURLAndUpdateStats(ctx context.Context, name string) (*url.URL, error)
}

func NewRepository(c *firestore.Client) Repository {
	return &repository{
		firestore: c,
	}
}

type repository struct {
	firestore *firestore.Client
}

type golink struct {
	URL string `firestore:"url"`
}

func (r *repository) GetURLAndUpdateStats(ctx context.Context, name string) (*url.URL, error) {
	col := r.firestore.Collection(collectionName)
	doc := col.Doc(name)

	s, err := doc.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, errors.Wrapf(errDocumentNotFound, "not found %s", doc.Path)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s", doc.Path)
	}

	var g golink
	if err := s.DataTo(&g); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", doc.Path)
	}

	u, err := url.Parse(g.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid url: %s", g.URL)
	}

	go r.incrementCount(ctx, doc)

	return u, nil
}

func (r *repository) incrementCount(ctx context.Context, docRef *firestore.DocumentRef) {
	_, err := docRef.Update(ctx, []firestore.Update{
		{
			Path:  "redirect_count",
			Value: firestore.Increment(1),
		},
	})
	if err != nil {
		err := errors.Wrapf(err, "failed to increment redirect_count of %s", docRef.Path)
		clog.Err(ctx, err)
	}
}

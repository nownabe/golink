package api

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/nownabe/golink/go/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collectionName = "golinks"

var errDocumentNotFound = errors.NewWithoutStack("not found")

type Repository interface {
	Transaction(ctx context.Context, f func(ctx context.Context, tx *firestore.Transaction) error) error
	Exists(ctx context.Context, tx *firestore.Transaction, name string) (bool, error)
	Get(ctx context.Context, name string) (*dto, error)
	Create(ctx context.Context, tx *firestore.Transaction, dto *dto) error
	ListByOwner(ctx context.Context, owner string) ([]*dto, error)
	ListByURL(ctx context.Context, tx *firestore.Transaction, url string) ([]*dto, error)
	Update(ctx context.Context, tx *firestore.Transaction, dto *dto) error
	Delete(ctx context.Context, tx *firestore.Transaction, name string) error
	AddOwner(ctx context.Context, tx *firestore.Transaction, name string, owner string) error
	RemoveOwner(ctx context.Context, tx *firestore.Transaction, name string, owner string) error
}

func NewRepository(c *firestore.Client) Repository {
	return &repository{
		firestore: c,
	}
}

type repository struct {
	firestore *firestore.Client
}

func (r *repository) Transaction(ctx context.Context, f func(ctx context.Context, tx *firestore.Transaction) error) error {
	return r.firestore.RunTransaction(ctx, f)
}

func (r *repository) Exists(ctx context.Context, tx *firestore.Transaction, name string) (bool, error) {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	s, err := tx.Get(doc)
	if status.Code(err) == codes.NotFound {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "failed to get golinks/%s", name)
	}

	return s.Exists(), nil
}

func (r *repository) Get(ctx context.Context, name string) (*dto, error) {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	s, err := doc.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, errors.Wrapf(errDocumentNotFound, "golinks/%s not found", name)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get golinks/%s", name)
	}

	var o dto
	if err := s.DataTo(&o); err != nil {
		return nil, errors.Wrapf(err, "failed to populate golinks/%s", name)
	}

	return &o, nil
}

func (r *repository) Create(ctx context.Context, tx *firestore.Transaction, dto *dto) error {
	col := r.collection()
	doc := col.Doc(dto.ID())

	dto.CreatedAt = time.Now()
	dto.UpdatedAt = time.Now()

	if err := tx.Create(doc, dto); err != nil {
		return errors.Wrapf(err, "failed to create golinks/%s", dto.Name)
	}

	return nil
}

func (r *repository) ListByOwner(ctx context.Context, owner string) ([]*dto, error) {
	col := r.collection()
	iter := col.Where("owners", "array-contains", owner).Documents(ctx)
	defer iter.Stop()

	var dtos []*dto
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to iterate golinks")
		}

		var o dto
		if err := s.DataTo(&o); err != nil {
			return nil, errors.Wrapf(err, "failed to populate golinks")
		}

		dtos = append(dtos, &o)
	}

	return dtos, nil
}

func (r *repository) ListByURL(ctx context.Context, tx *firestore.Transaction, url string) ([]*dto, error) {
	return nil, nil
}

func (r *repository) Update(ctx context.Context, tx *firestore.Transaction, dto *dto) error {
	return nil
}

func (r *repository) Delete(ctx context.Context, tx *firestore.Transaction, name string) error {
	return nil
}

func (r *repository) AddOwner(ctx context.Context, tx *firestore.Transaction, name string, owner string) error {
	return nil
}

func (r *repository) RemoveOwner(ctx context.Context, tx *firestore.Transaction, name string, owner string) error {
	return nil
}

func (r *repository) collection() *firestore.CollectionRef {
	return r.firestore.Collection(collectionName)
}

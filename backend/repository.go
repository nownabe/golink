package backend

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"go.nownabe.dev/clog/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionName                    = "golinks"
	firestoreFieldRedirectCount28Days = "redirect_count_28days"
	firestoreFieldRedirectCount7Days  = "redirect_count_7days"
)

var errDocumentNotFound = errors.NewWithoutStack("document not found")

type repository struct {
	firestore *firestore.Client
}

func (r *repository) Transaction(
	ctx context.Context,
	f func(ctx context.Context, tx *firestore.Transaction) error,
) error {
	if err := r.firestore.RunTransaction(ctx, f); err != nil {
		return errors.Errorf("r.firestore.RunTransaction(ctx, f): %w", err)
	}

	return nil
}

func (r *repository) TxExists(ctx context.Context, tx *firestore.Transaction, name string) (bool, error) {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	s, err := tx.Get(doc)
	if status.Code(err) == codes.NotFound {
		return false, nil
	}
	if err != nil {
		return false, errors.Errorf("failed to get %s: %w", doc.Path, err)
	}

	return s.Exists(), nil
}

func (r *repository) Get(ctx context.Context, name string) (*dto, error) {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	s, err := doc.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, errors.Errorf("%s not found: %w", doc.Path, errDocumentNotFound)
	}

	if err != nil {
		return nil, errors.Errorf("failed to get %s: %w", doc.Path, err)
	}
	var o dto
	if err := s.DataTo(&o); err != nil {
		return nil, errors.Errorf("failed to populate %s: %w", doc.Path, err)
	}

	return &o, nil
}

func (r *repository) TxGet(ctx context.Context, tx *firestore.Transaction, name string) (*dto, error) {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	s, err := tx.Get(doc)
	if status.Code(err) == codes.NotFound {
		return nil, errors.Errorf("%s not found: %w", doc.Path, errDocumentNotFound)
	}
	if err != nil {
		return nil, errors.Errorf("failed to get %s: %w", doc.Path, err)
	}

	var o dto
	if err := s.DataTo(&o); err != nil {
		return nil, errors.Errorf("failed to populate %s: %w", doc.Path, err)
	}

	return &o, nil
}

func (r *repository) TxCreate(ctx context.Context, tx *firestore.Transaction, dto *dto) error {
	col := r.collection()
	doc := col.Doc(dto.ID())

	dto.RedirectCountCalculatedDate = time.Now().Truncate(24 * time.Hour)
	dto.CreatedAt = time.Now()
	dto.UpdatedAt = time.Now()

	if err := tx.Create(doc, dto); err != nil {
		return errors.Errorf("failed to create %s: %w", doc.Path, err)
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
			return nil, errors.Errorf("failed to iterate %s: %w", col.Path, err)
		}

		var o dto
		if err := s.DataTo(&o); err != nil {
			return nil, errors.Errorf("failed to populate %s: %w", s.Ref.Path, err)
		}

		dtos = append(dtos, &o)
	}

	return dtos, nil
}

func (r *repository) ListByURL(ctx context.Context, url string) ([]*dto, error) {
	col := r.collection()
	iter := col.Where("url", "==", url).Documents(ctx)
	defer iter.Stop()

	var dtos []*dto
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Errorf("failed to iterate %s: %w", col.Path, err)
		}

		var o dto
		if err := s.DataTo(&o); err != nil {
			return nil, errors.Errorf("failed to populate %s: %w", s.Ref.Path, err)
		}

		dtos = append(dtos, &o)
	}

	return dtos, nil
}

func (r *repository) ListPopularGolinks(ctx context.Context, days, limit int) ([]*dto, error) {
	col := r.collection()
	var field string

	switch days {
	case 7:
		field = firestoreFieldRedirectCount7Days
	case 28:
		field = firestoreFieldRedirectCount28Days
	default:
		return nil, errors.Errorf("invalid days: %d", days)
	}

	iter := col.OrderBy(field, firestore.Desc).Limit(limit).Documents(ctx)
	defer iter.Stop()

	var golinks []*dto
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Errorf("failed to iterate %s: %w", col.Path, err)
		}

		var golink dto
		if err := s.DataTo(&golink); err != nil {
			return nil, errors.Errorf("failed to populate %s: %w", s.Ref.Path, err)
		}

		golinks = append(golinks, &golink)
	}

	return golinks, nil
}

func (r *repository) TxUpdate(ctx context.Context, tx *firestore.Transaction, dto *dto) error {
	col := r.collection()
	doc := col.Doc(dto.ID())

	dto.UpdatedAt = time.Now()

	if err := tx.Update(doc, []firestore.Update{
		{Path: "url", Value: dto.URL},
		{Path: firestoreFieldRedirectCount28Days, Value: dto.RedirectCount28Days},
		{Path: firestoreFieldRedirectCount7Days, Value: dto.RedirectCount7Days},
		{Path: "redirect_count_calculated_date", Value: dto.RedirectCountCalculatedDate},
		{Path: "daily_redirect_counts", Value: dto.DailyRedirectCounts},
		{Path: "updated_at", Value: dto.UpdatedAt},
	}); err != nil {
		return errors.Errorf("failed to update %s: %w", doc.Path, err)
	}

	return nil
}

func (r *repository) TxDelete(ctx context.Context, tx *firestore.Transaction, name string) error {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	if err := tx.Delete(doc); err != nil {
		return errors.Errorf("failed to delete %s: %w", doc.Path, err)
	}

	return nil
}

func (r *repository) TxAddOwner(ctx context.Context, tx *firestore.Transaction, name string, owner string) error {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	if err := tx.Update(doc, []firestore.Update{
		{Path: "owners", Value: firestore.ArrayUnion(owner)},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return errors.Errorf("failed to update %s: %w", doc.Path, err)
	}

	return nil
}

func (r *repository) TxRemoveOwner(ctx context.Context, tx *firestore.Transaction, name string, owner string) error {
	col := r.collection()
	doc := col.Doc(nameToID(name))

	if err := tx.Update(doc, []firestore.Update{
		{Path: "owners", Value: firestore.ArrayRemove(owner)},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return errors.Errorf("failed to update %s: %w", doc.Path, err)
	}

	return nil
}

func (r *repository) collection() *firestore.CollectionRef {
	return r.firestore.Collection(collectionName)
}

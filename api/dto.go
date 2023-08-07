package api

import (
	"time"

	golinkv1 "github.com/nownabe/golink/api/gen/golink/v1"
)

type dto struct {
	Name          string    `firestore:"-"`
	URL           string    `firestore:"url"`
	RedirectCount int64     `firestore:"redirect_count"`
	CreatedAt     time.Time `firestore:"created_at"`
	UpdatedAt     time.Time `firestore:"updated_at"`
	Owners        []string  `firestore:"owners"`
}

func (o *dto) toProto() *golinkv1.Golink {
	return &golinkv1.Golink{
		Name:   o.Name,
		Url:    o.URL,
		Owners: o.Owners,
	}
}

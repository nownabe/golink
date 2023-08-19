package backend

import (
	"strings"
	"time"

	golinkv1 "github.com/nownabe/golink/backend/gen/golink/v1"
)

type dto struct {
	Name          string    `firestore:"name"`
	URL           string    `firestore:"url"`
	RedirectCount int64     `firestore:"redirect_count"`
	CreatedAt     time.Time `firestore:"created_at"`
	UpdatedAt     time.Time `firestore:"updated_at"`
	Owners        []string  `firestore:"owners"`
}

func (o *dto) ID() string {
	return nameToID(o.Name)
}

func (o *dto) ToProto() *golinkv1.Golink {
	return &golinkv1.Golink{
		Name:   o.Name,
		Url:    o.URL,
		Owners: o.Owners,
	}
}

func nameToID(name string) string {
	return strings.ToLower(name)
}

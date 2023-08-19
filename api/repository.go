package api

import (
	"cloud.google.com/go/firestore"
	"github.com/nownabe/golink/go/errors"
)

const collectionName = "golinks"

var errDocumentNotFound = errors.NewWithoutStack("not found")

type repository struct {
	firestore *firestore.Client
}

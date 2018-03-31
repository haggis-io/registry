package service

import (
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/api"
)

type DocumentService interface {
	GetDocuments(dq model.DocumentQuery) ([]*api.Document, error)
	GetDocument(dq model.DocumentQuery) (*api.Document, error)
}

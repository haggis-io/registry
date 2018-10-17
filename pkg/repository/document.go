package repository

import (
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/model"
)

type DocumentRepository interface {
	Create(*api.Document) error
	Read(model.DocumentQuery) (*api.Document, error)
	Update(*api.Document) error
	Delete(model.DocumentQuery) error
	List(model.DocumentQuery) ([]*api.Document, error)
}

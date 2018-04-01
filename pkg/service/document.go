package service

import (
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/repository"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/haggis-io/registry/pkg/errors"
)

type registryService struct {
	db                 *gorm.DB
	documentRepository repository.CRUDL
}

func NewRegistryService(db *gorm.DB, documentRepository repository.CRUDL) DocumentService {
	return &registryService{
		db:                 db,
		documentRepository: documentRepository,
	}
}

func (r *registryService) GetDocuments(dq model.DocumentQuery) (out []*api.Document, err error) {
	query, err := dq.GenerateDetailedQuery()

	if err != nil {
		return out, status.Error(codes.InvalidArgument, err.Error())
	}

	documents, err := r.documentRepository.List(r.db, query)

	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return
	}

	return repository.ConvertSliceInterfaceToDocumentSlice(documents), nil

}

func (r *registryService) GetDocument(dq model.DocumentQuery) (out *api.Document, err error) {
	query, err := dq.GenerateDetailedQuery()

	if err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	document, err := r.documentRepository.Read(r.db, query)

	if err != nil {
		if err == errors.DocumentNotFoundErr {
			err = status.Error(codes.NotFound, err.Error())
			return
		}

		err = status.Error(codes.Internal, err.Error())
		return
	}

	return repository.ConvertDocumentToDocumentMessage(document.(*model.Document)), nil

}

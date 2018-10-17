package service

import (
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/errors"
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DocumentServiceInterface interface {
	List(dq model.DocumentQuery) ([]*api.Document, error)
	Get(dq model.DocumentQuery) (*api.Document, error)
}

func NewDocumentService(documentRepository repository.DocumentRepository) *DocumentService {
	return &DocumentService{
		documentRepository: documentRepository,
	}
}

type DocumentService struct {
	documentRepository repository.DocumentRepository
}

func (r *DocumentService) List(dq model.DocumentQuery) (out []*api.Document, err error) {
	return r.documentRepository.List(dq)
}

func (r *DocumentService) Get(dq model.DocumentQuery) (out *api.Document, err error) {

	out, err = r.documentRepository.Read(dq)

	if err != nil {
		if err == errors.DocumentNotFoundErr {
			err = status.Error(codes.NotFound, err.Error())
			return
		}

		err = status.Error(codes.Internal, err.Error())
		return
	}

	return
}

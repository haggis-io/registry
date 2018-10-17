package server

import (
	"context"

	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/service"
)

type RegistryServer struct {
	service service.DocumentServiceInterface
}

func NewRegistryServer(documentService service.DocumentServiceInterface) *RegistryServer {
	return &RegistryServer{
		service: documentService,
	}
}

func (r *RegistryServer) GetDocuments(ctx context.Context, request *api.GetDocumentsRequest) (*api.GetDocumentsResponse, error) {
	query := model.DocumentQuery{
		Name:    request.GetName(),
		Version: request.GetVersion(),
		Limit:   request.GetLimit(),
		Author:  request.GetAuthor(),
		Status:  request.GetStatus().String(),
	}

	docs, err := r.service.List(query)

	if err != nil {
		return nil, err
	}

	return &api.GetDocumentsResponse{
		Documents: docs,
	}, nil

}
func (r *RegistryServer) GetDocument(ctx context.Context, request *api.GetDocumentRequest) (*api.GetDocumentResponse, error) {
	query := model.DocumentQuery{
		Name:    request.GetName(),
		Version: request.GetVersion(),
	}

	doc, err := r.service.Get(query)

	if err != nil {
		return nil, err
	}

	return &api.GetDocumentResponse{
		Document: doc,
	}, nil
}

func (r *RegistryServer) CreateDocument(ctx context.Context, request *api.CreateDocumentRequest) (*api.GetDocumentResponse, error) {
	return nil, nil
}

func (r *RegistryServer) MarkAsApproved(ctx context.Context, request *api.ApprovedDocumentRequest) (*api.ApprovedDocumentResponse, error) {
	return nil, nil
}

func (r *RegistryServer) MarkAsDeclined(ctx context.Context, request *api.DeclinedDocumentRequest) (*api.DeclinedDocumentResponse, error) {
	return nil, nil
}

func (r *RegistryServer) MarkAsPending(ctx context.Context, request *api.PendingDocumentRequest) (*api.PendingDocumentResponse, error) {
	return nil, nil
}

package server

import (
	"context"
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/service"
)

type registryServer struct {
	service service.DocumentService
}

func NewRegistryServer(documentService service.DocumentService) api.RegistryServer {
	return &registryServer{
		service: documentService,
	}
}

func (r *registryServer) GetDocuments(ctx context.Context, request *api.GetDocumentsRequest) (*api.GetDocumentsResponse, error) {
	query := model.DocumentQuery{
		Name:    request.GetName(),
		Version: model.Version(request.GetVersion()),
		Limit:   int64(request.GetLimit()),
		Author:  request.GetAuthor(),
		Status:  model.StatusMap[int32(request.GetStatus())],
	}

	docs, err := r.service.GetDocuments(query)

	if err != nil {
		return nil, err
	}

	return &api.GetDocumentsResponse{
		Documents: docs,
	}, nil

}
func (r *registryServer) GetDocument(ctx context.Context, request *api.GetDocumentRequest) (*api.GetDocumentResponse, error) {
	return nil, nil
}
func (r *registryServer) CreateDocument(ctx context.Context, request *api.CreateDocumentRequest) (*api.GetDocumentResponse, error) {
	return nil, nil
}
func (r *registryServer) MarkAsApproved(ctx context.Context, request *api.ApprovedDocumentRequest) (*api.ApprovedDocumentResponse, error) {
	return nil, nil
}
func (r *registryServer) MarkAsDeclined(ctx context.Context, request *api.DeclinedDocumentRequest) (*api.DeclinedDocumentResponse, error) {
	return nil, nil
}
func (r *registryServer) MarkAsPending(ctx context.Context, request *api.PendingDocumentRequest) (*api.PendingDocumentResponse, error) {
	return nil, nil
}

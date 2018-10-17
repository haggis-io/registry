package relational

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/storage/relational/entity"
)

func documentToDocumentMessage(document *entity.Document) *api.Document {

	createdAt, err := ptypes.TimestampProto(document.CreatedAt)

	if err != nil {
		return nil
	}

	updatedAt, err := ptypes.TimestampProto(document.UpdatedAt)

	if err != nil {
		return nil
	}

	out := &api.Document{
		Name:        document.Name,
		Description: document.Description,
		Version:     string(document.Version),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Author:      document.Author,
		Status:      api.Status(api.Status_value[string(document.Status)]),
		Snippet: &api.Snippet{
			Text:     document.Snippet.Text,
			TestCase: document.Snippet.TestCase,
		},
	}

	for _, dep := range document.Dependencies {
		out.Dependencies = append(out.Dependencies, documentToDocumentMessage(dep))
	}

	return out

}

func documentMessageToDocument(document *api.Document) *entity.Document {

	createdAt, err := ptypes.Timestamp(document.CreatedAt)

	if err != nil {
		return nil
	}

	updatedAt, err := ptypes.Timestamp(document.UpdatedAt)

	if err != nil {
		return nil
	}

	out := &entity.Document{
		Name:        document.Name,
		Description: document.Description,
		Version:     entity.Version(document.Version),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Author:      document.Author,
		Status:      entity.Status(document.Status),
		Snippet: entity.Snippet{
			Text:     document.Snippet.Text,
			TestCase: document.Snippet.TestCase,
		},
	}

	for _, dep := range document.Dependencies {
		out.Dependencies = append(out.Dependencies, documentMessageToDocument(dep))
	}

	return out

}

func documentsTodocumentMessages(documents []*entity.Document) []*api.Document {
	out := make([]*api.Document, len(documents))

	for i, document := range documents {
		out[i] = documentToDocumentMessage(document)
	}

	return out
}

func documentMessagesTodocuments(documents []*api.Document) []*entity.Document {
	out := make([]*entity.Document, len(documents))

	for i, document := range documents {
		out[i] = documentMessageToDocument(document)
	}

	return out
}

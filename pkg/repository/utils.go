package repository

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/model"
)

func ConvertDocumentSliceToSliceInterface(docs []*model.Document) []interface{} {
	out := make([]interface{}, len(docs))

	for i, v := range docs {
		out[i] = v
	}

	return out
}

func ConvertSliceInterfaceToDocumentSlice(interfaces []interface{}) []*api.Document {
	out := make([]*api.Document, len(interfaces))

	for i, v := range interfaces {
		document := v.(*model.Document)

		createdAt, err := ptypes.TimestampProto(document.CreatedAt)

		if err != nil {
			return out
		}

		updatedAt, err := ptypes.TimestampProto(document.UpdatedAt)

		if err != nil {
			return out
		}

		out[i] = &api.Document{
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
	}

	return out
}

func ConvertDocumentToDocumentMessage(document *model.Document) *api.Document {

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
		out.Dependencies = append(out.Dependencies, ConvertDocumentToDocumentMessage(dep))
	}

	return out

}

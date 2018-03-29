package document

import (
	"github.com/haggis-io/registry/pkg/model"
	"github.com/jinzhu/gorm"
)

type DocumentRepository struct{}

func (d *DocumentRepository) Create(tx *gorm.DB, document *model.Document) error {
	return tx.Create(document).Error
}

func (d *DocumentRepository) GetDocuments(tx *gorm.DB, name string) ([]*model.Document, error) {
	var out []*model.Document
	err := tx.Preload("Dependencies").Find(&out, "name = ?", name).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return out, nil
		}
		return nil, err
	}
	return out, nil
}

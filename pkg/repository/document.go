package repository

import (
	"github.com/haggis-io/registry/pkg/model"
	"github.com/jinzhu/gorm"
)

type documentRepository struct{}

func NewDocumentRepository() CRUDL {
	return &documentRepository{}
}

func (*documentRepository) Create(tx *gorm.DB, document interface{}) error {
	return tx.Create(document).Error
}

func (*documentRepository) Read(tx *gorm.DB, query model.Query) (interface{}, error) {
	var out *model.Document

	err := query(tx).Find(&out).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return out, nil
		}
		return nil, err
	}
	return out, nil
}

func (*documentRepository) Update(tx *gorm.DB, document interface{}) error {
	return tx.Save(document).Error
}

func (*documentRepository) Delete(tx *gorm.DB, query model.Query) error {
	return query(tx).Delete(model.Document{}).Error
}

func (*documentRepository) List(tx *gorm.DB, query model.Query) ([]interface{}, error) {
	var out []*model.Document

	err := query(tx).Find(&out).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return ConvertDocumentSliceToSliceInterface(out), nil
}

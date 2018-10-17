package relational

import (
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/errors"
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/storage/relational/entity"
	"github.com/jinzhu/gorm"
)

type Query func(*gorm.DB) *gorm.DB

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (d *DocumentRepository) Create(doc *api.Document) error {
	return d.db.Create(doc).Error
}

func (d *DocumentRepository) Read(query model.DocumentQuery) (*api.Document, error) {
	var out entity.Document

	detailedQuery, err := generateDetailedQuery(query)

	if err != nil {
		return nil, err
	}

	err = detailedQuery(d.db).Find(&out).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return documentToDocumentMessage(&out), errors.DocumentNotFoundErr
		}
		return nil, err
	}
	return documentToDocumentMessage(&out), nil
}

func (d *DocumentRepository) Update(doc *api.Document) error {
	return d.db.Save(doc).Error
}

func (d *DocumentRepository) Delete(query model.DocumentQuery) error {
	detailedQuery, err := generateDetailedQuery(query)

	if err != nil {
		return err
	}
	return detailedQuery(d.db).Delete(entity.Document{}).Error
}

func (d *DocumentRepository) List(query model.DocumentQuery) ([]*api.Document, error) {
	var out []*entity.Document

	detailedQuery, err := generateDetailedQuery(query)

	if err != nil {
		return nil, err
	}

	err = detailedQuery(d.db).Find(&out).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.DocumentNotFoundErr
		}
		return nil, err
	}
	return documentsTodocumentMessages(out), nil
}

func generateDetailedQuery(d model.DocumentQuery) (Query, error) {

	var (
		query   string
		where   []interface{}
		version entity.Version = entity.Version(d.Version)
		status  entity.Status  = entity.Status(d.Status)
	)

	if d.Name != "" {

		if d.Pattern {
			query = "name LIKE ? AND"
			where = append(where, "%"+d.Name+"%")
		} else {
			query = "name = ? AND"
			where = append(where, d.Name)
		}

		if version != "" {

			if err := version.Validate(); err != nil {
				return nil, err
			}
			query += " version = ? AND"
			where = append(where, version)
		}
	}

	if d.Author != "" {
		query += " author = ? AND"
		where = append(where, d.Author)
	}

	query += " status = ?"

	if status != "" {
		where = append(where, status)
	} else {
		where = append(where, entity.StatusAPPROVED)
	}

	if d.Limit > 0 {
		return func(tx *gorm.DB) *gorm.DB {
			return addRequiredPreloads(tx).Limit(d.Limit).Where(query, where...)
		}, nil
	} else if len(where) > 0 {
		return func(tx *gorm.DB) *gorm.DB {
			return addRequiredPreloads(tx).Where(query, where...)
		}, nil
	}

	return func(tx *gorm.DB) *gorm.DB {
		return addRequiredPreloads(tx).Where(query, where...)
	}, nil

}

func addRequiredPreloads(tx *gorm.DB) *gorm.DB {
	var (
		snippetPreloadField = "Snippet"
		dependencyDepth1    = "Dependencies.Snippet"
		dependencyDepth2    = "Dependencies.Dependencies.Snippet"
		dependencyDepth3    = "Dependencies.Dependencies.Dependencies.Snippet"
		dependencyDepth4    = "Dependencies.Dependencies.Dependencies.Dependencies.Snippet"
		dependencyDepth5    = "Dependencies.Dependencies.Dependencies.Dependencies.Dependencies.Snippet"
	)
	return tx.
		Preload(snippetPreloadField).
		Preload(dependencyDepth1).
		Preload(dependencyDepth2).
		Preload(dependencyDepth3).
		Preload(dependencyDepth4).
		Preload(dependencyDepth5)

}

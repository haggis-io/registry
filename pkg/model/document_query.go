package model

import "github.com/jinzhu/gorm"

type DocumentQuery struct {
	Name    string
	Pattern bool
	Version Version
	Limit   int64
	Author  string
	Status  Status
}

func (d DocumentQuery) GenerateDetailedQuery() (Query, error) {

	var (
		query string
		where []interface{}
	)

	if d.Name != "" {

		if d.Pattern {
			query = "name LIKE ? AND"
			where = append(where, "%"+d.Name+"%")
		} else {
			query = "name = ? AND"
			where = append(where, d.Name)
		}

		if d.Version != "" {
			if err := d.Version.Validate(); err != nil {
				return nil, err
			}
			query += " version = ? AND"
			where = append(where, d.Version)
		}
	}

	if d.Author != "" {
		query += " author = ? AND"
		where = append(where, d.Author)
	}

	query += " status = ?"

	if d.Status != "" {
		where = append(where, d.Status)
	} else {
		where = append(where, StatusAPPROVED)
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
		SnippetField    = "Snippet"
		dependencyField = "Dependencies.Snippet"
	)

	return tx.Preload(SnippetField).Preload(dependencyField)

}

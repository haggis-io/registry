package repository

import (
	"github.com/haggis-io/registry/pkg/model"
	"github.com/jinzhu/gorm"
)

type CRUDL interface {
	Create(*gorm.DB, interface{}) error
	Read(*gorm.DB, model.Query) (interface{}, error)
	Update(*gorm.DB, interface{}) error
	Delete(*gorm.DB, model.Query) error
	List(*gorm.DB, model.Query) ([]interface{}, error)
}

package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type Document struct {
	Name string `gorm:"primary_key"`

	Description string

	Version Version `gorm:"primary_key"`

	Author string

	Status Status

	Dependencies []*Document `gorm:"many2many:document_dependencies;association_jointable_foreignkey:dependency_name,dependency_version;association_autoupdate:false;association_autocreate:false"`

	Snippet Snippet `gorm:"foreignkey:DocumentName,DocumentVersion"`

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt *time.Time

	Helper map[string]interface{} `gorm:"-"` //ignored
}

func (d *Document) BeforeCreate(scope *gorm.Scope) error {
	d.Status.Pending()
	scope.SetColumn("NAME", d.Name)
	scope.SetColumn("VERSION", d.Version)
	return nil
}

func (d *Document) BeforeSave(scope *gorm.Scope) error {
	s, ok := d.Helper["Version"].(string)

	if !ok {
		return errors.New("wrong type assertion")
	}

	d.Version = Version(s)

	if err := d.Version.Validate(); err != nil {
		return err
	}

	d.Name = d.Helper["Name"].(string)

	return nil
}

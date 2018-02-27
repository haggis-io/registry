package repository

import (
	"errors"
	"github.com/haggis-io/registry/pkg/proto"
	"upper.io/db.v3/lib/sqlbuilder"
)

var (
	TableNotExistError = errors.New("table does not exist")
)

type EntityRepository interface {
	GetEntities(entityType, author string, status registry.Status, limit int32) ([]*registry.Entity, error)
	GetEntityByName(entityType, name string) (*registry.Entity, error)
	CreateEntity(entity *registry.Entity) (*registry.Entity, error)
	DeleteEntity(entity *registry.Entity) error
}

type entityRepository struct {
	db sqlbuilder.Database
}

func NewEntityRepository(db sqlbuilder.Database) EntityRepository {
	return &entityRepository{
		db: db,
	}
}

func (r *entityRepository) GetEntities(entityType, author string, status registry.Status, limit int32) ([]*registry.Entity, error) {
	var (
		result     []*registry.Entity
		entityColl = r.db.Collection(entityType)
	)

	if !entityColl.Exists() {
		return nil, TableNotExistError
	}

	res := entityColl.Find("author = ? and status = ?", author, status)

	if err := res.Limit(int(limit)).All(&result); err != nil {
		return nil, TableNotExistError
	}

	return result, nil
}

func (r *entityRepository) GetEntityByName(entityType, name string) (*registry.Entity, error) {
	var (
		result     *registry.Entity
		entityColl = r.db.Collection(entityType)
	)

	if !entityColl.Exists() {
		return nil, TableNotExistError
	}

	res := entityColl.Find("name", name)

	if err := res.One(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *entityRepository) CreateEntity(entity *registry.Entity) (*registry.Entity, error) {
	var (
		entityColl = r.db.Collection(entity.GetType())
	)

	if !entityColl.Exists() {
		return nil, TableNotExistError
	}

	if err := entityColl.InsertReturning(entity); err != nil {
		return nil, err
	}

	return entity, nil
}
func (r *entityRepository) DeleteEntity(entity *registry.Entity) error {
	var (
		entityColl = r.db.Collection(entity.GetType())
	)

	if !entityColl.Exists() {
		return TableNotExistError
	}

	res := entityColl.Find("name", entity.GetName())

	return res.Delete()
}

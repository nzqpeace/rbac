package db

import (
	"math"

	"github.com/nzqpeace/rbac/rbac/model"
	"gopkg.in/mgo.v2/bson"
)

const PermissionsList = "permissions"

type PermissionDao struct {
	*Base
}

func NewPermissionDao(db *DataBase) *PermissionDao {
	return &PermissionDao{
		NewBase(db, PermissionsList),
	}
}

func (dao *PermissionDao) GetAllPermissions(system string) (ps []model.Permission, err error) {
	err = dao.FindAll(bson.M{"system": system}, &ps, 0, math.MaxInt32)
	return
}

func (dao *PermissionDao) CreatePermission(p *model.Permission) error {
	return dao.Upsert(bson.M{"system": p.System, "name": p.Name}, p)
}

func (dao *PermissionDao) RemovePermission(system, name string) error {
	return dao.Remove(bson.M{"system": system, "name": name})
}

func (dao *PermissionDao) UpdatePermission(system, oldname, newname string) error {
	return dao.Update(bson.M{"system": system, "name": oldname}, bson.M{
		"$set": bson.M{
			"name": newname,
		},
	})
}

package db

import (
	"math"

	"github.com/nzqpeace/rbac/rbac/model"
	"gopkg.in/mgo.v2/bson"
)

const RoleList = "roles"

type RoleDao struct {
	*Base
}

func NewRoleDao(db *DataBase) *RoleDao {
	return &RoleDao{
		NewBase(db, RoleList),
	}
}

func (dao *RoleDao) GetRole(system, name string) (role model.Role, err error) {
	err = dao.Find(bson.M{"system": system, "name": name}, &role)
	return
}

func (dao *RoleDao) GetAllRoles(system string) (roles []model.Role, err error) {
	err = dao.FindAll(bson.M{"system": system}, &roles, 0, math.MaxInt32)
	return
}

func (dao *RoleDao) CreateRole(role *model.Role) error {
	return dao.Upsert(bson.M{"system": role.System, "name": role.Name}, role)
}

func (dao *RoleDao) RemoveRole(system, name string) error {
	return dao.Remove(bson.M{"system": system, "name": name})
}

func (dao *RoleDao) RemoveAllRoles(system string) error {
	return dao.RemoveAll(bson.M{"system": system})
}

func (dao *RoleDao) UpdateRoleName(system, oldname, newname string) error {
	return dao.Update(bson.M{"system": system, "name": oldname},
		bson.M{
			"$set": bson.M{
				"name": newname,
			},
		},
	)
}

func (dao *RoleDao) GetPermissions(system, name string) ([]string, error) {
	var role model.Role
	err := dao.Find(bson.M{"system": system, "name": name}, &role)
	return role.Permissions, err
}

func (dao *RoleDao) GrantPermissions(system, name string, permissions ...string) error {
	if len(permissions) == 0 {
		return nil
	}

	return dao.Update(bson.M{"system": system, "name": name}, bson.M{
		"$push": bson.M{
			"permissions": bson.M{
				"$each": permissions,
			},
		},
	})
}

func (dao *RoleDao) RemovePermission(system, name string, permission string) error {
	return dao.Update(bson.M{"system": system, "name": name}, bson.M{
		"$pull": bson.M{
			"permissions": permission,
		},
	})
}

package cache

import (
	"fmt"

	set "github.com/deckarep/golang-set"
	"github.com/nzqpeace/rbac/rbac/db"
	"github.com/nzqpeace/rbac/rbac/model"
)

const redisKeyFormatPermissions = "%s_%s_permissions" // {system}_{uid}_permissions

// PermissionDao is permission dao
type PermissionDao struct {
	*Redis
	role *db.RoleDao
	user *db.UserDao
}

// NewPermissionDao create a new permission dao
func NewPermissionDao(r *Redis, mgo *db.DataBase) *PermissionDao {
	return &PermissionDao{
		r,
		db.NewRoleDao(mgo),
		db.NewUserDao(mgo),
	}
}

// Permissions list all permissions of specified system
func (dao *PermissionDao) Permissions(system, uid string) (ps []string, err error) {
	key := fmt.Sprintf(redisKeyFormatPermissions, system, uid)
	return dao.SMembers(key)
}

// IsPermit check if have specified permission
func (dao *PermissionDao) IsPermit(system, uid string, permission string) (permit bool, err error) {
	key := fmt.Sprintf(redisKeyFormatPermissions, system, uid)
	permit, err = dao.SIsMembers(key, permission)
	if err != nil {
		return
	}
	// check whether specified key exist when `exist` is false
	if !permit {
		exist, err := dao.Exists(key)
		if err != nil {
			return permit, err
		}

		if !exist { // reload from mongo when specified key is not in cache
			err = dao.ReloadPermissions(system, uid)
			return dao.SIsMembers(key, permission)
		}
	}
	return
}

// RemovePermissions remove specified permissions
func (dao *PermissionDao) RemovePermissions(system, uid string, names ...string) error {
	key := fmt.Sprintf(redisKeyFormatPermissions, system, uid)
	return dao.SRem(key, names...)
}

// AddPermissions add specified permissions
func (dao *PermissionDao) AddPermissions(system, uid string, names ...string) error {
	key := fmt.Sprintf(redisKeyFormatPermissions, system, uid)
	return dao.SAdd(key, names...)
}

// ReloadPermissions reload permissions from mongo
func (dao *PermissionDao) ReloadPermissions(system, uid string) error {
	key := fmt.Sprintf(redisKeyFormatPermissions, system, uid)
	_, err := dao.Del(key)
	if err != nil {
		return err
	}

	// reload from mongo
	userPermModel, err := dao.user.GetUserPermModel(system, uid)
	if err != nil {
		return err
	}

	permissions := dao.GetPermissions(&userPermModel)

	// store permissions into redis
	return dao.SAdd(key, permissions...)
}

func (dao *PermissionDao) GetPermissions(u *model.UserPermModel) (permissions []string) {
	pset := set.NewSet()
	// generate permission list
	// 1. add permissions at whitelist
	for _, p := range u.WhiteList {
		pset.Add(p)
	}

	// 2. add permissions permited throught role
	for _, role := range u.Roles {
		// fetch each role's permissions
		ps, err := dao.role.GetPermissions(u.System, role)
		if err != nil {
			continue
		}

		for _, p := range ps {
			pset.Add(p)
		}
	}

	// 3. remove permissions at blacklist
	for _, p := range u.BlackList {
		pset.Remove(p)
	}

	for _, v := range pset.ToSlice() {
		permissions = append(permissions, v.(string))
	}
	return
}

func (dao *PermissionDao) RemoveUser(system, uid string) (bool, error) {
	key := fmt.Sprintf(redisKeyFormatPermissions, system, uid)
	return dao.Del(key)
}

func (dao *PermissionDao) ClearAllKeys() {
	dao.FlushDB()
}

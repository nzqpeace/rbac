package db

import (
	"github.com/nzqpeace/rbac/model"
	"gopkg.in/mgo.v2/bson"
)

// UserList is name of collection
const UserList = "user"

// UserDao define dao of user
type UserDao struct {
	*Base
}

// NewUserDao create a new instance of UserDao
func NewUserDao(db *DataBase) *UserDao {
	return &UserDao{
		NewBase(db, UserList),
	}
}

// CreateUserPermModel associate user info with permission, and store it into db
func (dao *UserDao) CreateUserPermModel(user *model.UserPermModel) error {
	return dao.Upsert(bson.M{"system": user.System, "uid": user.UID}, user)
}

// RemoveUserPermModel remove user info from mongo
func (dao *UserDao) RemoveUserPermModel(system, uid string) error {
	return dao.Remove(bson.M{"system": system, "uid": uid})
}

// UpdateUserPermModel update user info
func (dao *UserDao) UpdateUserPermModel(system, uid string, user *model.UserPermModel) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, user)
}

// GetUserPermModel get user info
func (dao *UserDao) GetUserPermModel(system, uid string) (user model.UserPermModel, err error) {
	err = dao.Find(bson.M{"system": system, "uid": uid}, &user)
	return
}

// GetAllRoles get all roles with uid
func (dao *UserDao) GetAllRoles(system, uid string) (roles []string, err error) {
	var user model.UserPermModel
	err = dao.Find(bson.M{"system": system, "uid": uid}, &user)
	if err == nil {
		roles = user.Roles
	}
	return
}

// UpdateRoles update user's all roles
func (dao *UserDao) UpdateRoles(system, uid string, roles ...string) error {
	if len(roles) == 0 {
		return nil
	}

	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$set": bson.M{
			"roles": roles,
		},
	})
}

// AddRoles add specified roles into user's permission model
func (dao *UserDao) AddRoles(system, uid string, roles ...string) error {
	if len(roles) == 0 {
		return nil
	}

	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$push": bson.M{
			"roles": bson.M{
				"$each": roles,
			},
		},
	})
}

// RemoveRoles remove specified role from user's permission model
func (dao *UserDao) RemoveRoles(system, uid string, role string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$pull": bson.M{
			"roles": role,
		},
	})
}

// GetBlackList get user permission model's blacklist, which contain all permissions forbidden
func (dao *UserDao) GetBlackList(system, uid string) (bl []string, err error) {
	var user model.UserPermModel
	err = dao.Find(bson.M{"system": system, "uid": uid}, &user)
	if err == nil {
		bl = user.BlackList
	}
	return
}

// AddToBlackList add specified permissions into user permission model's blacklist
func (dao *UserDao) AddToBlackList(system, uid string, permissions ...string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$push": bson.M{
			"blacklist": bson.M{
				"$each": permissions,
			},
		},
	})
}

// RemoveFromBlackList remove specified permission from blacklist
func (dao *UserDao) RemoveFromBlackList(system, uid string, permission string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$pull": bson.M{
			"blacklist": permission,
		},
	})
}

// ClearBlackList clear blacklist
func (dao *UserDao) ClearBlackList(system, uid string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$set": bson.M{
			"blacklist": []string{},
		},
	})
}

// GetWhiteList get user permission model's whitelist, which contain all permissions allowed all the time
func (dao *UserDao) GetWhiteList(system, uid string) (whitelist []string, err error) {
	var user model.UserPermModel
	err = dao.Find(bson.M{"system": system, "uid": uid}, &user)
	if err == nil {
		whitelist = user.WhiteList
	}
	return
}

// UpdateWhiteList update whitelist with 'wl'
func (dao *UserDao) UpdateWhiteList(system, uid string, whitelist ...string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$set": bson.M{
			"whitelist": whitelist,
		},
	})
}

// AddToWhiteList add specified permission into user permission model's whitelist
func (dao *UserDao) AddToWhiteList(system, uid string, permissions ...string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$push": bson.M{
			"whitelist": bson.M{
				"$each": permissions,
			},
		},
	})
}

// RemoveFromWhiteList remove specified permission from user's permission model's whitelist
func (dao *UserDao) RemoveFromWhiteList(system, uid string, permission string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$pull": bson.M{
			"whitelist": permission,
		},
	})
}

// ClearWhiteList clear all permission at user's permission model's whitelist
func (dao *UserDao) ClearWhiteList(system, uid string) error {
	return dao.Update(bson.M{"system": system, "uid": uid}, bson.M{
		"$set": bson.M{
			"whitelist": []string{},
		},
	})
}

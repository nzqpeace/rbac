package db

import (
	"fmt"
	"testing"

	"github.com/nzqpeace/rbac/model"
	"github.com/stretchr/testify/assert"
)

var (
	userDao *UserDao
)

func init() {
	conf = &MgoConf{
		Url: "localhost/test",
	}

	var err error
	db, err = Init(conf)
	if err != nil {
		fmt.Println(err)
	}

	pdao = NewPermissionDao(db)
	roleDao = NewRoleDao(db)
	userDao = NewUserDao(db)
}

func fillUserData(t *testing.T) {
	guestUser := model.NewUserPermModel(system, "uid_guest", "guest")
	commonUser := model.NewUserPermModel(system, "uid_common", "common")
	adminUser := model.NewUserPermModel(system, "uid_admin", "common", "admin")

	// create user permission model
	assert.Nil(t, userDao.CreateUserPermModel(guestUser))
	assert.Nil(t, userDao.CreateUserPermModel(commonUser))
	assert.Nil(t, userDao.CreateUserPermModel(adminUser))
}

func clearUserData(t *testing.T) {
	assert.Nil(t, userDao.RemoveUserPermModel(system, "uid_guest"))
	assert.Nil(t, userDao.RemoveUserPermModel(system, "uid_common"))
	assert.Nil(t, userDao.RemoveUserPermModel(system, "uid_admin"))
}

func TestUserPermModel(t *testing.T) {
	fillUserData(t)

	// get user permission model
	u, err := userDao.GetUserPermModel(system, "uid_guest")
	assert.Nil(t, err)
	assert.Equal(t, "uid_guest", u.UID)

	// update roles
	assert.Nil(t, userDao.UpdateRoles(system, "uid_common", "manage"))

	roles, err := userDao.GetAllRoles(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(roles))
	assert.Equal(t, "manage", roles[0])

	// add roles
	assert.Nil(t, userDao.AddRoles(system, "uid_common", "write", "manage"))
	roles, err = userDao.GetAllRoles(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(roles))

	// remove roles
	assert.Nil(t, userDao.RemoveRoles(system, "uid_common", "manage"))
	roles, err = userDao.GetAllRoles(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(roles))

	// add to blacklist
	assert.Nil(t, userDao.AddToBlackList(system, "uid_common", "write"))
	bl, err := userDao.GetBlackList(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bl))
	assert.Equal(t, "write", bl[0])

	// remove from blacklist
	assert.Nil(t, userDao.RemoveFromBlackList(system, "uid_common", "write"))
	bl, err = userDao.GetBlackList(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bl))

	// add to whitelist
	assert.Nil(t, userDao.AddToWhiteList(system, "uid_common", "write"))
	wl, err := userDao.GetWhiteList(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(wl))
	assert.Equal(t, "write", wl[0])

	// remove from whitelist
	assert.Nil(t, userDao.RemoveFromWhiteList(system, "uid_common", "write"))
	wl, err = userDao.GetWhiteList(system, "uid_common")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(wl))

	clearUserData(t)
}

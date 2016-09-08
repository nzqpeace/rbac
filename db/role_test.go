package db

import (
	"fmt"
	"testing"

	"github.com/nzqpeace/rbac/model"
	"github.com/stretchr/testify/assert"
)

var (
	roleDao *RoleDao
	guest   *model.Role
	common  *model.Role
	admin   *model.Role
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

	roleDao = NewRoleDao(db)
}

func fillRoleData(t *testing.T) {
	// define roles
	guest = model.NewRole(system, "guest", "", "read")
	common = model.NewRole(system, "common", "", "read", "write")
	admin = model.NewRole(system, "admin", "", "read", "write", "manage")

	// create role
	assert.Nil(t, roleDao.CreateRole(guest))
	assert.Nil(t, roleDao.CreateRole(common))
	assert.Nil(t, roleDao.CreateRole(admin))
}

func clearRoleData(t *testing.T) {
	// remove role
	assert.Nil(t, roleDao.RemoveRole(system, "guest"))
	assert.Nil(t, roleDao.RemoveAllRoles(system))
}

func TestRole(t *testing.T) {
	fillRoleData(t)

	// query role
	r, err := roleDao.GetRole(system, "common")
	assert.Nil(t, err)
	assert.Equal(t, common.Name, r.Name)
	assert.Equal(t, common.Desc, r.Desc)
	assert.Equal(t, 2, len(r.Permissions))

	// get all roles
	rs, err := roleDao.GetAllRoles(system)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(rs))

	// update role
	assert.Nil(t, roleDao.UpdateRoleName(system, "common", "user"))
	r, err = roleDao.GetRole(system, "common")
	assert.NotNil(t, err)

	r, err = roleDao.GetRole(system, "user")
	assert.Nil(t, err)
	assert.NotEqual(t, guest.Name, r.Name)
	assert.Equal(t, "user", r.Name)

	// grant permissions
	assert.Nil(t, roleDao.GrantPermissions(system, "guest", "write", "manage"))
	r, err = roleDao.GetRole(system, "guest")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(r.Permissions))

	// remove permission
	assert.Nil(t, roleDao.RemovePermission(system, "guest", "manage"))
	r, err = roleDao.GetRole(system, "guest")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(r.Permissions))

	clearRoleData(t)
}

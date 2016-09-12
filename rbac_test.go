package rbac

import (
	"fmt"
	"testing"

	"github.com/nzqpeace/rbac/cache"
	"github.com/nzqpeace/rbac/db"
	"github.com/stretchr/testify/assert"
)

const (
	system = "Cowshed"

	// permissions
	read   = "read"
	write  = "write"
	manage = "manage"

	// roles
	guest  = "guest"
	common = "common"
	admin  = "admin"

	// uid
	uid_guest  = "uid_guest"
	uid_common = "uid_common"
	uid_admin  = "uid_admin"
)

var (
	rbac *RBAC
)

func init() {
	conf := &RBACConfig{
		Redis: cache.DefaultConfig(),
		Mgo: &db.MgoConf{
			Url: "localhost/test",
		},
	}

	var err error
	rbac, err = NewRBAC(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func fillTestData(t *testing.T) {
	assert.NotNil(t, rbac)
	// register permissions
	assert.Nil(t, rbac.RegisterPermission(system, read, "read question/answer/comment"))
	assert.Nil(t, rbac.RegisterPermission(system, write, "post question/answer/comment"))
	assert.Nil(t, rbac.RegisterPermission(system, manage, "manage question and answer"))

	// register roles
	assert.Nil(t, rbac.RegisterRole(system, guest, "", read))
	assert.Nil(t, rbac.RegisterRole(system, common, "", read, write))
	assert.Nil(t, rbac.RegisterRole(system, admin, "", read, write, manage))

	// register users
	assert.Nil(t, rbac.RegisterUser(system, uid_guest, guest))
	assert.Nil(t, rbac.RegisterUser(system, uid_common, common))
	assert.Nil(t, rbac.RegisterUser(system, uid_admin, common, admin))
}

func clearTestData(t *testing.T) {
	// remove all permissions
	assert.Nil(t, rbac.UnregisterPermission(system, read))
	assert.Nil(t, rbac.UnregisterPermission(system, write))
	assert.Nil(t, rbac.UnregisterPermission(system, manage))

	// remove all roles
	assert.Nil(t, rbac.UnregisterRole(system, guest))
	assert.Nil(t, rbac.UnregisterAllRoles(system))

	// remove all users
	assert.Nil(t, rbac.UnregisterUser(system, uid_guest))
	assert.Nil(t, rbac.UnregisterUser(system, uid_common))
	assert.Nil(t, rbac.UnregisterUser(system, uid_admin))
}

func TestRBAC(t *testing.T) {
	fillTestData(t)

	// IsPermit
	permit, err := rbac.IsPermit(system, uid_common, manage)
	assert.Nil(t, err)
	assert.False(t, permit)

	permit, err = rbac.IsPermit(system, uid_common, write)
	assert.Nil(t, err)
	assert.True(t, permit)

	// GetAllPermissionsBySystem
	p, err := rbac.GetAllPermissionsBySystem(system)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(p))

	// GetAllRolesOfSystem
	role, err := rbac.GetRoleOfSystem(system, admin)
	assert.Nil(t, err)
	assert.Equal(t, admin, role.Name)

	// GetPermissionsOfRole
	ps, err := rbac.GetPermissionsOfRole(system, admin)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ps))

	// GrantPermissionsToRole
	assert.Nil(t, rbac.GrantPermissionsToRole(system, guest, write))
	ps, err = rbac.GetPermissionsOfRole(system, guest)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ps))

	permit, err = rbac.IsPermit(system, uid_guest, write)
	assert.Nil(t, err)
	assert.True(t, permit)

	// RemovePermissionFromRole
	assert.Nil(t, rbac.RemovePermissionFromRole(system, guest, write))
	ps, err = rbac.GetPermissionsOfRole(system, guest)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ps))

	permit, err = rbac.IsPermit(system, uid_guest, write)
	assert.Nil(t, err)
	assert.False(t, permit)

	// GetAllRolesByUID
	rs, err := rbac.GetAllRolesByUID(system, uid_admin)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rs))

	// AddRoles
	assert.Nil(t, rbac.AddRoles(system, uid_guest, common))
	rs, err = rbac.GetAllRolesByUID(system, uid_guest)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rs))

	permit, err = rbac.IsPermit(system, uid_guest, write)
	assert.Nil(t, err)
	assert.True(t, permit)

	// RemoveRoles
	assert.Nil(t, rbac.RemoveRoles(system, uid_guest, common))
	rs, err = rbac.GetAllRolesByUID(system, uid_guest)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(rs))

	permit, err = rbac.IsPermit(system, uid_guest, write)
	assert.Nil(t, err)
	assert.False(t, permit)

	// AddToBlackList
	assert.Nil(t, rbac.AddToBlackList(system, uid_admin, write, manage))

	permit, err = rbac.IsPermit(system, uid_admin, write)
	assert.Nil(t, err)
	assert.False(t, permit)

	permit, err = rbac.IsPermit(system, uid_admin, manage)
	assert.Nil(t, err)
	assert.False(t, permit)

	// GetBlackList
	bl, err := rbac.GetBlackList(system, uid_admin)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(bl))

	// RemoveFromBlackList
	assert.Nil(t, rbac.RemoveFromBlackList(system, uid_admin, manage))
	bl, err = rbac.GetBlackList(system, uid_admin)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bl))

	permit, err = rbac.IsPermit(system, uid_admin, manage)
	assert.Nil(t, err)
	assert.True(t, permit)

	// ClearBlackList
	assert.Nil(t, rbac.ClearBlackList(system, uid_admin))
	bl, err = rbac.GetBlackList(system, uid_admin)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bl))

	permit, err = rbac.IsPermit(system, uid_admin, manage)
	assert.Nil(t, err)
	assert.True(t, permit)

	// AddToWhiteList
	assert.Nil(t, rbac.AddToWhiteList(system, uid_guest, write, manage))
	permit, err = rbac.IsPermit(system, uid_guest, write)
	assert.Nil(t, err)
	assert.True(t, permit)

	permit, err = rbac.IsPermit(system, uid_guest, manage)
	assert.Nil(t, err)
	assert.True(t, permit)

	// GetWhiteList
	wl, err := rbac.GetWhiteList(system, uid_guest)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(wl))

	// RemoveFromWhiteList
	assert.Nil(t, rbac.RemoveFromWhiteList(system, uid_guest, manage))
	wl, err = rbac.GetWhiteList(system, uid_guest)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(wl))

	permit, err = rbac.IsPermit(system, uid_guest, manage)
	assert.Nil(t, err)
	assert.False(t, permit)

	// ClearWhiteList
	assert.Nil(t, rbac.ClearWhiteList(system, uid_guest))
	wl, err = rbac.GetWhiteList(system, uid_guest)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(wl))

	permit, err = rbac.IsPermit(system, uid_guest, write)
	assert.Nil(t, err)
	assert.False(t, permit)

	clearTestData(t)
}

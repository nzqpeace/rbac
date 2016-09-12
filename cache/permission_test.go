package cache

import (
	"fmt"
	"testing"

	"github.com/nzqpeace/rbac/db"
	"github.com/nzqpeace/rbac/model"
	"github.com/stretchr/testify/assert"
)

var (
	pdao *PermissionDao
)

const system = "Cowshed"
const uid = "uid_common"

func init() {
	config := &RedisConfig{
		Address:       ":6379",
		Password:      "",
		DB:            1,
		MaxConn:       100,
		IdleTimeout:   60,
		RetryInterval: 3,
		RetryTimes:    0,
	}

	r := NewRedis(config)

	conf := &db.MgoConf{
		Url: "localhost/test",
	}

	db, err := db.Init(conf)
	if err != nil {
		fmt.Println(err)
	}

	pdao = NewPermissionDao(r, db)
}

func fillDataIntoMongo(t *testing.T) {
	// define roles
	guest := model.NewRole(system, "guest", "", "read")
	common := model.NewRole(system, "common", "", "read", "write")
	admin := model.NewRole(system, "admin", "", "read", "write", "manage")

	// create role
	assert.Nil(t, pdao.role.CreateRole(guest))
	assert.Nil(t, pdao.role.CreateRole(common))
	assert.Nil(t, pdao.role.CreateRole(admin))

	commonUser := model.NewUserPermModel(system, uid, "common")
	commonAdmin := model.NewUserPermModel(system, "uid_admin", "admin")

	// create user permission model
	assert.Nil(t, pdao.user.CreateUserPermModel(commonUser))
	assert.Nil(t, pdao.user.CreateUserPermModel(commonAdmin))
}

func clearDataAtMongo(t *testing.T) {
	assert.Nil(t, pdao.user.RemoveUserPermModel(system, uid))
	assert.Nil(t, pdao.role.RemoveAllRoles(system))
}

func TestPermission(t *testing.T) {
	// add permissions
	assert.Nil(t, pdao.AddPermissions(system, uid, "read", "write"))

	// remove permissions
	assert.Nil(t, pdao.RemovePermissions(system, uid, "write"))

	// get all permissions
	ps, err := pdao.Permissions(system, uid)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ps))
	assert.Equal(t, "read", ps[0])

	// is permited
	permit, err := pdao.IsPermit(system, uid, "read")
	assert.Nil(t, err)
	assert.True(t, permit)

	permit, err = pdao.IsPermit(system, uid, "write")
	assert.Nil(t, err)
	assert.False(t, permit)

	// test reload
	fillDataIntoMongo(t)

	assert.Nil(t, pdao.ReloadPermissions(system, uid))
	// get all permissions
	ps, err = pdao.Permissions(system, uid)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ps))

	// is permited
	permit, err = pdao.IsPermit(system, uid, "read")
	assert.Nil(t, err)
	assert.True(t, permit)

	permit, err = pdao.IsPermit(system, uid, "write")
	assert.Nil(t, err)
	assert.True(t, permit)

	permit, err = pdao.IsPermit(system, uid, "manage")
	assert.Nil(t, err)
	assert.False(t, permit)

	// key is not exist
	permit, err = pdao.IsPermit(system, "uid_not_exist", "manage")
	assert.Nil(t, err)
	assert.False(t, permit)

	permit, err = pdao.IsPermit(system, "uid_admin", "manage")
	assert.Nil(t, err)
	assert.True(t, permit)

	clearDataAtMongo(t)

	// remove permissions
	assert.Nil(t, pdao.RemovePermissions(system, uid, "write", "read"))
}

func BenchmarkIsPermit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// pdao.SIsMembers("cowshed_uid_admin_permissions", "read")
		pdao.IsPermit("cowshed", "uid_admin", "read")
	}
}

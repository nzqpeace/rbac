package db

import (
	"fmt"
	"testing"

	"github.com/nzqpeace/rbac/model"
	"github.com/stretchr/testify/assert"
)

var (
	pdao *PermissionDao
)

const system = "Cowshed"

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
}

func fillPermissionData(t *testing.T) {
	read := &model.Permission{
		System: system,
		Name:   "read",
		Desc:   "read question/answer/comment",
	}

	write := &model.Permission{
		System: system,
		Name:   "write",
		Desc:   "post question/answer/comment",
	}

	manage := &model.Permission{
		System: system,
		Name:   "manage",
		Desc:   "manage question and answer",
	}

	// insert test
	assert.Nil(t, pdao.CreatePermission(read))
	assert.Nil(t, pdao.CreatePermission(write))
	assert.Nil(t, pdao.CreatePermission(manage))
}

func clearPermissionData(t *testing.T) {
	// remove all documents
	assert.Nil(t, pdao.RemovePermission(system, "read"))
	assert.Nil(t, pdao.RemovePermission(system, "write"))
	assert.Nil(t, pdao.RemovePermission(system, "admin"))
}

func TestPermission(t *testing.T) {
	fillPermissionData(t)

	// update
	assert.Nil(t, pdao.UpdatePermission(system, "manage", "admin"))
	assert.NotNil(t, pdao.RemovePermission(system, "manage"))

	// get all permissions
	ps, err := pdao.GetAllPermissions(system)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(ps))

	clearPermissionData(t)
}

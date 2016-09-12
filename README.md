# RBAC  [![Build Status](https://travis-ci.org/nzqpeace/rbac.svg?branch=master)](https://travis-ci.org/nzqpeace/rbac)

Role-Based Access Control (RBAC) for Golang

# Usage
```Golang
package main

import (
	"fmt"

	"github.com/nzqpeace/rbac"
	"github.com/nzqpeace/rbac/cache"
	"github.com/nzqpeace/rbac/db"
)

var (
	r *rbac.RBAC
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

func init() {
	conf := &rbac.RBACConfig{
		Redis: cache.DefaultConfig(),
		Mgo: &db.MgoConf{
			Url: "localhost/test",
		},
	}

	var err error
	r, err = rbac.NewRBAC(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// register permissions
	r.RegisterPermission(system, read, "read question/answer/comment")
	r.RegisterPermission(system, write, "post question/answer/comment")
	r.RegisterPermission(system, manage, "manage question and answer")

	// register roles
	r.RegisterRole(system, guest, "", read)
	r.RegisterRole(system, common, "", read, write)
	r.RegisterRole(system, admin, "", read, write, manage)

	// register users
	r.RegisterUser(system, uid_guest, guest)
	r.RegisterUser(system, uid_common, common)
	r.RegisterUser(system, uid_admin, common, admin)

	// check
	permit, err := r.IsPermit(system, uid_common, manage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(permit) // false

	permit, err = r.IsPermit(system, uid_common, write)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(permit) // true

	permit, err = r.IsPermit(system, uid_common, read)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(permit) // true
}

```
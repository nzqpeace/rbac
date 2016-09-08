package main

import (
	"fmt"

	"github.com/nzqpeace/devgroup/rbac"
)

type ForumPerm int

const (
	ForumPermRead ForumPerm = iota
	ForumPermPost
	ForumPermReply
	ForumPermComment
	ForumPermDelete
)

func (f ForumPerm) String() string {
	switch f {
	case ForumPermRead:
		return "read"
	case ForumPermPost:
		return "post"
	case ForumPermReply:
		return "reply"
	case ForumPermComment:
		return "comment"
	case ForumPermDelete:
		return "delete"
	default:
		return ""
	}
}

func main() {
	r := rbac.NewRBAC()

	permRead := rbac.NewPermission(int(ForumPermRead), ForumPermRead.String(), ForumPermRead.String()+"'s desc")
	permPost := rbac.NewPermission(int(ForumPermPost), ForumPermPost.String(), ForumPermPost.String()+"'s desc")
	permReply := rbac.NewPermission(int(ForumPermReply), ForumPermReply.String(), ForumPermReply.String()+"'s desc")
	permComment := rbac.NewPermission(int(ForumPermComment), ForumPermComment.String(), ForumPermComment.String()+"'s desc")
	permDelete := rbac.NewPermission(int(ForumPermDelete), ForumPermDelete.String(), ForumPermDelete.String()+"'s desc")

	roleGuest := rbac.NewRole("guest", "guest only can read post", permRead)
	roleCommon := rbac.NewRole("common", "common user can post and reply", permRead, permPost, permReply, permComment)
	roleAdmin := rbac.NewRole("admin", "admin has all permisions", permRead, permPost, permReply, permComment, permDelete)

	// register permissions and roles
	r.RegisterPermission(permRead)
	r.RegisterPermission(permPost)
	r.RegisterPermission(permReply)
	r.RegisterPermission(permComment)
	r.RegisterPermission(permDelete)

	r.RegisterRole(roleGuest)
	r.RegisterRole(roleCommon)
	r.RegisterRole(roleAdmin)

	// associated uid with rbac model
	u1 := rbac.NewUserPermModel("u1", roleGuest)
	u2 := rbac.NewUserPermModel("u2", roleCommon)
	u3 := rbac.NewUserPermModel("u3", roleAdmin)

	fmt.Println(u1.Permit(permPost)) // false
	fmt.Println(u1.Permit(permRead)) // true

	fmt.Println(u2.Permit(permPost)) // true
	u2.BL.Add(permPost)
	fmt.Println(u2.Permit(permPost)) // false
	u2.BL.Remove(permPost)
	fmt.Println(u2.Permit(permPost)) // true

	u3.WL.Add(permDelete)
}

package rbac

import (
	"github.com/nzqpeace/rbac/cache"
	"github.com/nzqpeace/rbac/db"
	"github.com/nzqpeace/rbac/model"
)

// RBAC entry of rbac system
type RBAC struct {
	Cache      *cache.PermissionDao
	Permission *db.PermissionDao
	Role       *db.RoleDao
	User       *db.UserDao
}

// NewRBAC create a new instance
func NewRBAC(config *RBACConfig) (rbac *RBAC, err error) {
	r := cache.NewRedis(config.Redis)
	d, err := db.Init(config.Mgo)
	if err != nil {
		return nil, err
	}

	rbac = &RBAC{
		Cache:      cache.NewPermissionDao(r, d),
		Permission: db.NewPermissionDao(d),
		Role:       db.NewRoleDao(d),
		User:       db.NewUserDao(d),
	}
	return
}

// IsPermit check whether have specified permission
func (r *RBAC) IsPermit(system, uid, permission string) (bool, error) {
	return r.Cache.IsPermit(system, uid, permission)
}

// RegisterPermission register permission
func (r *RBAC) RegisterPermission(system, name, desc string) error {
	p := &model.Permission{
		System: system,
		Name:   name,
		Desc:   desc,
	}
	return r.Permission.CreatePermission(p)
}

// UnregisterPermission remove permission from system
func (r *RBAC) UnregisterPermission(system, permission string) error {
	return r.Permission.RemovePermission(system, permission)
}

// GetAllPermissionsBySystem get all permissions of specified system
func (r *RBAC) GetAllPermissionsBySystem(system string) ([]model.Permission, error) {
	return r.Permission.GetAllPermissions(system)
}

// UpdatePermission update permission
func (r *RBAC) UpdatePermission(system, oldname, newname string) error {
	return r.Permission.UpdatePermission(system, oldname, newname)
}

// RegisterRole register role
func (r *RBAC) RegisterRole(system, name, desc string, permissions ...string) error {
	role := model.NewRole(system, name, desc, permissions...)
	return r.Role.CreateRole(role)
}

// UnregisterRole unregister specified role of specified system
func (r *RBAC) UnregisterRole(system, name string) error {
	return r.Role.RemoveRole(system, name)
}

// UnregisterAllRoles unregister all role of specified system
func (r *RBAC) UnregisterAllRoles(system string) error {
	return r.Role.RemoveAllRoles(system)
}

// GetRoleOfSystem get specified role of system by name
func (r *RBAC) GetRoleOfSystem(system, name string) (model.Role, error) {
	return r.Role.GetRole(system, name)
}

// GetAllRolesOfSystem get all roles of specified system
func (r *RBAC) GetAllRolesOfSystem(system string) ([]model.Role, error) {
	return r.Role.GetAllRoles(system)
}

// UpdateRoleName update name of specified role
func (r *RBAC) UpdateRoleName(system, oldname, newname string) error {
	return r.Role.UpdateRoleName(system, oldname, newname)
}

// GetPermissionsOfRole get all permissions of role
func (r *RBAC) GetPermissionsOfRole(system, name string) ([]string, error) {
	return r.Role.GetPermissions(system, name)
}

// GrantPermissionsToRole grant specified permissions to role
func (r *RBAC) GrantPermissionsToRole(system, name string, permissions ...string) error {
	r.Cache.ClearAllKeys()
	return r.Role.GrantPermissions(system, name, permissions...)
}

// RemovePermissionFromRole remove specified permission from specified role
func (r *RBAC) RemovePermissionFromRole(system, name string, permission string) error {
	r.Cache.ClearAllKeys()
	return r.Role.RemovePermission(system, name, permission)
}

// RegisterUser register user permission info into mongo
func (r *RBAC) RegisterUser(system, uid string, roles ...string) error {
	u := model.NewUserPermModel(system, uid, roles...)
	return r.User.CreateUserPermModel(u)
}

// UnregisterUser remove user info from mongo
func (r *RBAC) UnregisterUser(system, uid string) error {
	return r.User.RemoveUserPermModel(system, uid)
}

// UpdateUser update user info
func (r *RBAC) UpdateUser(system, uid string, new_roles ...string) error {
	u := model.NewUserPermModel(system, uid, new_roles...)
	return r.User.UpdateUserPermModel(system, uid, u)
}

// GetUser get user info
func (r *RBAC) GetUser(system, uid string) (model.UserPermModel, error) {
	return r.User.GetUserPermModel(system, uid)
}

// GetAllRolesByUID get all roles with uid
func (r *RBAC) GetAllRolesByUID(system, uid string) (roles []string, err error) {
	return r.User.GetAllRoles(system, uid)
}

// UpdateRoles update user's all roles
func (r *RBAC) UpdateRoles(system, uid string, roles ...string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.UpdateRoles(system, uid, roles...)
}

// AddRoles add specified roles into user's permission model
func (r *RBAC) AddRoles(system, uid string, roles ...string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.AddRoles(system, uid, roles...)
}

// RemoveRoles remove specified role from user's permission model
func (r *RBAC) RemoveRoles(system, uid string, role string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.RemoveRoles(system, uid, role)
}

// GetBlackList get user permission model's blacklist, which contain all permissions forbidden
func (r *RBAC) GetBlackList(system, uid string) ([]string, error) {
	return r.User.GetBlackList(system, uid)
}

// AddToBlackList add specified permissions into user permission model's blacklist
func (r *RBAC) AddToBlackList(system, uid string, permissions ...string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.AddToBlackList(system, uid, permissions...)
}

// RemoveFromBlackList remove specified permission from blacklist
func (r *RBAC) RemoveFromBlackList(system, uid string, permission string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.RemoveFromBlackList(system, uid, permission)
}

// ClearBlackList clear blacklist
func (r *RBAC) ClearBlackList(system, uid string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.ClearBlackList(system, uid)
}

// GetWhiteList get user permission model's whitelist, which contain all permissions allowed all the time
func (r *RBAC) GetWhiteList(system, uid string) ([]string, error) {
	return r.User.GetWhiteList(system, uid)
}

// UpdateWhiteList update whitelist with 'wl'
func (r *RBAC) UpdateWhiteList(system, uid string, whitelist ...string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.UpdateWhiteList(system, uid, whitelist...)
}

// AddToWhiteList add specified permission into user permission model's whitelist
func (r *RBAC) AddToWhiteList(system, uid string, permissions ...string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.AddToWhiteList(system, uid, permissions...)
}

// RemoveFromWhiteList remove specified permission from user's permission model's whitelist
func (r *RBAC) RemoveFromWhiteList(system, uid string, permission string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.RemoveFromWhiteList(system, uid, permission)
}

// ClearWhiteList clear all permission at user's permission model's whitelist
func (r *RBAC) ClearWhiteList(system, uid string) error {
	r.Cache.RemoveUser(system, uid)
	return r.User.ClearWhiteList(system, uid)
}

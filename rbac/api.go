package main

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/nzqpeace/rbac"
	"github.com/nzqpeace/rbac/model"

	"github.com/kataras/iris"
)

type RbacApi struct {
	rbac *rbac.RBAC
}

type ErrCode int

const (
	NotFound  = "not found"
	Success   = "success"
	BadParams = "bad params"
)

const (
	ErrOK ErrCode = iota
	ErrNotFound
	ErrBadPrams
	ErrInternelServerError
)

func NewRbacApi(config *Config) (*RbacApi, error) {
	rc := &rbac.RBACConfig{
		Redis: config.Redis,
		Mgo:   config.Mongo,
	}

	r, err := rbac.NewRBAC(rc)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &RbacApi{r}, nil
}

func validateParams(c iris.Context, params interface{}) error {
	if err := c.ReadJSON(params); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{
			"code":    ErrBadPrams,
			"message": err.Error(),
		})
		return err
	}

	if err := validate.Struct(params); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{
			"code":    ErrBadPrams,
			"message": err.Error(),
		})
		return err
	}

	return nil
}

func checkUrlParams(c iris.Context, names ...string) (params map[string]string, err error) {
	params = make(map[string]string)
	for _, name := range names {
		value := c.URLParam(name)
		if value == "" {
			message := fmt.Sprintf("miss parameter[%s]", name)
			c.StatusCode(iris.StatusBadRequest)
			c.JSON(iris.Map{
				"code":    ErrBadPrams,
				"message": message,
			})
			return params, errors.New(message)
		}
		params[name] = value
	}
	return
}

func (api *RbacApi) responseByError(c iris.Context, err error) {
	if err != nil {
		if strings.Contains(err.Error(), NotFound) {
			c.StatusCode(iris.StatusOK)
			c.JSON(iris.Map{
				"code":    ErrNotFound,
				"message": err.Error(),
			})
			return
		}
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{
			"code":    ErrInternelServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(iris.Map{
		"code":    ErrOK,
		"message": Success,
	})
}

func (api *RbacApi) responseAdditionData(c iris.Context, err error, jsonKey string, jsonValue interface{}) {
	if err != nil {
		if strings.Contains(err.Error(), NotFound) {
			c.StatusCode(iris.StatusOK)
			c.JSON(iris.Map{
				"code":    ErrNotFound,
				"message": err.Error(),
			})
			return
		}

		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{
			"code":    ErrInternelServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(iris.Map{
		"code":    ErrOK,
		"message": Success,
		jsonKey:   jsonValue,
	})
}

// IsPermit check whether have specified permission
func (api *RbacApi) IsPermit(c iris.Context) {
	params, err := checkUrlParams(c, "system", "uid", "permission")
	if err != nil {
		return
	}

	permit, err := api.rbac.IsPermit(params["system"], params["uid"], params["permission"])
	api.responseAdditionData(c, err, "permit", permit)
}

// RegisterPermission register permission
func (api *RbacApi) RegisterPermission(c iris.Context) {
	var p model.Permission
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.RegisterPermission(p.System, p.Name, p.Desc)
	api.responseByError(c, err)
}

// UnregisterPermission remove permission from system
func (api *RbacApi) UnregisterPermission(c iris.Context) {
	var p model.Permission
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UnregisterPermission(p.System, p.Name)
	api.responseByError(c, err)
}

// GetAllPermissionsBySystem get all permissions of specified system
func (api *RbacApi) GetAllPermissionsBySystem(c iris.Context) {
	params, err := checkUrlParams(c, "system")
	if err != nil {
		return
	}

	ps, err := api.rbac.GetAllPermissionsBySystem(params["system"])
	if err != nil && strings.Contains(err.Error(), NotFound) {
		c.StatusCode(iris.StatusOK)
		c.JSON(iris.Map{
			"code":    ErrNotFound,
			"message": err.Error(),
		})
		return
	}

	api.responseAdditionData(c, err, "permissions", ps)
}

// UpdatePermission update permission
func (api *RbacApi) UpdatePermission(c iris.Context) {
	var p struct {
		System  string `json:"system" validate:"required"`
		OldName string `json:"oldname" validate:"required"`
		NewName string `json:"newname" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UpdatePermission(p.System, p.OldName, p.NewName)
	api.responseByError(c, err)
}

// RegisterRole register role
func (api *RbacApi) RegisterRole(c iris.Context) {
	var role model.Role
	if validateParams(c, &role) != nil {
		return
	}

	err := api.rbac.RegisterRole(role.System, role.Name, role.Desc, role.Permissions...)
	api.responseByError(c, err)
}

// UnregisterRole unregister specified role of specified system
func (api *RbacApi) UnregisterRole(c iris.Context) {
	var role struct {
		System string `json:"system" validate:"required"`
		Name   string `json:"name" validate:"required"`
	}
	if validateParams(c, &role) != nil {
		return
	}

	err := api.rbac.UnregisterRole(role.System, role.Name)
	api.responseByError(c, err)
}

// UnregisterAllRoles unregister all role of specified system
func (api *RbacApi) UnregisterAllRoles(c iris.Context) {
	var role struct {
		System string `json:"system" validate:"required"`
		Name   string `json:"name" validate:"required"`
	}
	if validateParams(c, &role) != nil {
		return
	}

	err := api.rbac.UnregisterAllRoles(role.System)
	api.responseByError(c, err)
}

// GetRoleOfSystem get specified role of system by name
func (api *RbacApi) GetRoleOfSystem(c iris.Context) {
	params, err := checkUrlParams(c, "system", "role")
	if err != nil {
		return
	}

	role, err := api.rbac.GetRoleOfSystem(params["system"], params["role"])
	if err != nil && strings.Contains(err.Error(), NotFound) {
		c.StatusCode(iris.StatusOK)
		c.JSON(iris.Map{
			"code":    ErrNotFound,
			"message": err.Error(),
		})
		return
	}

	api.responseAdditionData(c, err, "role", role)
}

// GetAllRolesOfSystem get all roles of specified system
func (api *RbacApi) GetAllRolesOfSystem(c iris.Context) {
	params, err := checkUrlParams(c, "system")
	if err != nil {
		return
	}

	roles, err := api.rbac.GetAllRolesOfSystem(params["system"])
	if err != nil && strings.Contains(err.Error(), NotFound) {
		c.StatusCode(iris.StatusOK)
		c.JSON(iris.Map{
			"code":    ErrNotFound,
			"message": err.Error(),
		})
		return
	}

	api.responseAdditionData(c, err, "roles", roles)
}

// UpdateRoleName update name of specified role
func (api *RbacApi) UpdateRoleName(c iris.Context) {
	var p struct {
		System  string `json:"system" validate:"required"`
		OldName string `json:"oldname" validate:"required"`
		NewName string `json:"newname" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UpdateRoleName(p.System, p.OldName, p.NewName)
	api.responseByError(c, err)
}

// GetPermissionsOfRole get all permissions of role
func (api *RbacApi) GetPermissionsOfRole(c iris.Context) {
	params, err := checkUrlParams(c, "system", "role")
	if err != nil {
		return
	}

	ps, err := api.rbac.GetPermissionsOfRole(params["system"], params["role"])
	if err != nil && strings.Contains(err.Error(), NotFound) {
		c.StatusCode(iris.StatusOK)
		c.JSON(iris.Map{
			"code":    ErrNotFound,
			"message": err.Error(),
		})
		return
	}

	api.responseAdditionData(c, err, "permissions", ps)
}

// GrantPermissionsToRole grant specified permissions to role
func (api *RbacApi) GrantPermissionsToRole(c iris.Context) {
	var p struct {
		System      string   `json:"system" validate:"required"`
		Role        string   `json:"role" validate:"required"`
		Permissions []string `json:"permissions" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.GrantPermissionsToRole(p.System, p.Role, p.Permissions...)
	api.responseByError(c, err)
}

// RemovePermissionFromRole remove specified permission from specified role
func (api *RbacApi) RemovePermissionFromRole(c iris.Context) {
	var p struct {
		System     string `json:"system" validate:"required"`
		Role       string `json:"role" validate:"required"`
		Permission string `json:"permission" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.RemovePermissionFromRole(p.System, p.Role, p.Permission)
	api.responseByError(c, err)
}

// RegisterUser register user permission info into mongo
func (api *RbacApi) RegisterUser(c iris.Context) {
	var p model.UserPermModel
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.RegisterUser(p.System, p.UID, p.Roles...)
	api.responseByError(c, err)
}

// UnregisterUser remove user info from mongo
func (api *RbacApi) UnregisterUser(c iris.Context) {
	var p struct {
		System string `json:"system" validate:"required"`
		UID    string `json:"uid" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UnregisterUser(p.System, p.UID)
	api.responseByError(c, err)
}

// UpdateUser update user info
func (api *RbacApi) UpdateUser(c iris.Context) {
	var p struct {
		System   string   `json:"system" validate:"required"`
		UID      string   `json:"uid" validate:"required"`
		NewRoles []string `json:"new_roles" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UpdateUser(p.System, p.UID, p.NewRoles...)
	api.responseByError(c, err)
}

// GetUser get user info
func (api *RbacApi) GetUser(c iris.Context) {
	params, err := checkUrlParams(c, "system", "uid")
	if err != nil {
		return
	}

	u, err := api.rbac.GetUser(params["system"], params["uid"])
	api.responseAdditionData(c, err, "user", u)
}

// GetAllRolesByUID get all roles with uid
func (api *RbacApi) GetAllRolesByUID(c iris.Context) {
	params, err := checkUrlParams(c, "system", "uid")
	if err != nil {
		return
	}

	roles, err := api.rbac.GetAllRolesByUID(params["system"], params["uid"])
	api.responseAdditionData(c, err, "roles", roles)
}

// UpdateRoles update user's all roles
func (api *RbacApi) UpdateRoles(c iris.Context) {
	var p model.UserPermModel
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UpdateRoles(p.System, p.UID, p.Roles...)
	api.responseByError(c, err)
}

// AddRoles add specified roles into user's permission model
func (api *RbacApi) AddRoles(c iris.Context) {
	var p model.UserPermModel
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.AddRoles(p.System, p.UID, p.Roles...)
	api.responseByError(c, err)
}

// RemoveRoles remove specified role from user's permission model
func (api *RbacApi) RemoveRoles(c iris.Context) {
	var p struct {
		System string `json:"system"  validate:"required"`
		UID    string `json:"uid" validate:"required"`
		Role   string `json:"role" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.RemoveRoles(p.System, p.UID, p.Role)
	api.responseByError(c, err)
}

// GetBlackList get user permission model's blacklist, which contain all permissions forbidden
func (api *RbacApi) GetBlackList(c iris.Context) {
	params, err := checkUrlParams(c, "system", "uid")
	if err != nil {
		return
	}

	bl, err := api.rbac.GetBlackList(params["system"], params["uid"])
	api.responseAdditionData(c, err, "blacklist", bl)
}

// AddToBlackList add specified permissions into user permission model's blacklist
func (api *RbacApi) AddToBlackList(c iris.Context) {
	var p struct {
		System      string   `json:"system"  validate:"required"`
		UID         string   `json:"uid" validate:"required"`
		Permissions []string `json:"permissions"validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.AddToBlackList(p.System, p.UID, p.Permissions...)
	api.responseByError(c, err)
}

// RemoveFromBlackList remove specified permission from blacklist
func (api *RbacApi) RemoveFromBlackList(c iris.Context) {
	var p struct {
		System     string `json:"system"  validate:"required"`
		UID        string `json:"uid" validate:"required"`
		Permission string `json:"permission"validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.RemoveFromBlackList(p.System, p.UID, p.Permission)
	api.responseByError(c, err)
}

// ClearBlackList clear blacklist
func (api *RbacApi) ClearBlackList(c iris.Context) {
	var p struct {
		System string `json:"system"  validate:"required"`
		UID    string `json:"uid" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.ClearBlackList(p.System, p.UID)
	api.responseByError(c, err)
}

// GetWhiteList get user permission model's whitelist, which contain all permissions allowed all the time
func (api *RbacApi) GetWhiteList(c iris.Context) {
	params, err := checkUrlParams(c, "system", "uid")
	if err != nil {
		return
	}

	wl, err := api.rbac.GetWhiteList(params["system"], params["uid"])
	api.responseAdditionData(c, err, "whitelist", wl)
}

// UpdateWhiteList update whitelist with 'wl'
func (api *RbacApi) UpdateWhiteList(c iris.Context) {
	var p struct {
		System    string   `json:"system"  validate:"required"`
		UID       string   `json:"uid" validate:"required"`
		WhiteList []string `json:"whitelist"validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.UpdateWhiteList(p.System, p.UID, p.WhiteList...)
	api.responseByError(c, err)
}

// AddToWhiteList add specified permission into user permission model's whitelist
func (api *RbacApi) AddToWhiteList(c iris.Context) {
	var p struct {
		System      string   `json:"system"  validate:"required"`
		UID         string   `json:"uid" validate:"required"`
		Permissions []string `json:"permissions"validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.AddToWhiteList(p.System, p.UID, p.Permissions...)
	api.responseByError(c, err)
}

// RemoveFromWhiteList remove specified permission from user's permission model's whitelist
func (api *RbacApi) RemoveFromWhiteList(c iris.Context) {
	var p struct {
		System     string `json:"system"  validate:"required"`
		UID        string `json:"uid" validate:"required"`
		Permission string `json:"permission"validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.RemoveFromWhiteList(p.System, p.UID, p.Permission)
	api.responseByError(c, err)
}

// ClearWhiteList clear all permission at user's permission model's whitelist
func (api *RbacApi) ClearWhiteList(c iris.Context) {
	var p struct {
		System string `json:"system"  validate:"required"`
		UID    string `json:"uid" validate:"required"`
	}
	if validateParams(c, &p) != nil {
		return
	}

	err := api.rbac.ClearWhiteList(p.System, p.UID)
	api.responseByError(c, err)
}

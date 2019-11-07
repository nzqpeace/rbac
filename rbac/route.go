package main

import (
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

func registerRoute(app *iris.Application, config *Config) error {
	rbacAPI, err := NewRbacApi(config)
	if err != nil {
		log.Error(err)
		return err
	}

	// check check whether have specified permission
	// URL params: system, uid, permission
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "permit":true // true or false
	// }
	app.Get("/authenticate", rbacAPI.IsPermit)

	// register permission
	// Json params:
	// {
	//     "system":system,
	//     "name":name,
	//     "desc":description {option}
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Post("/permission", rbacAPI.RegisterPermission)

	// unregister permission
	// Json params:
	// {
	//     "system":system,
	//     "name":name
	//     "desc":description {option}
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Delete("/permission", rbacAPI.UnregisterPermission)

	// get all permissions by system
	// URL params: system
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "permissions":[
	//         {
	//             "system":system,
	//             "name":name
	//             "desc":desc
	//         },
	//         {
	//             "system":system,
	//             "name":name
	//             "desc":desc
	//         }
	//     ]
	// }
	app.Get("/permission", rbacAPI.GetAllPermissionsBySystem)

	// update permission
	// Json params
	// {
	//     "system":system,
	//     "oldname":oldname,
	//     "newname":newname
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/permission", rbacAPI.UpdatePermission)

	// register role
	// Json params:
	// {
	//     "system":system,
	//     "name":name,
	//     "desc":description {option}
	//     "permissions":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Post("/role", rbacAPI.RegisterRole)

	// unregister role
	// Json params:
	// {
	//     "system":system,
	//     "name":name
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Delete("/role", rbacAPI.UnregisterRole)

	// get specified role by system and role name
	// URL params: system, role
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "role":{
	//         "system":system,
	//         "name":name,
	//         "desc":desc,
	//         "permissions":[
	//             "permission1",
	//             "permission2"
	//         ]
	//     }
	// }
	app.Get("/role", rbacAPI.GetRoleOfSystem)

	// get all roles by system
	// URL params: system
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "roles":[
	//         {
	//             "system":system,
	//             "name":name,
	//             "desc":desc,
	//             "permissions":[
	//                 "permission1",
	//                 "permission2"
	//             ]
	//         },
	//         {
	//             "system":system,
	//             "name":name,
	//             "desc":desc,
	//             "permissions":[
	//                 "permission1",
	//                 "permission2"
	//             ]
	//         }
	//     ]
	// }
	app.Get("/role/all", rbacAPI.GetAllRolesOfSystem)

	// update role name
	// Json params:
	// {
	//     "system":system,
	//     "oldname":oldname,
	//     "newname":newname,
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/role", rbacAPI.UpdateRoleName)

	// get all permissions of specified role
	// URL params: system, role
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "permissions":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	app.Get("/role/permissions", rbacAPI.GetPermissionsOfRole)

	// grant permissions to specified role
	// Json params:
	// {
	//     "system":system,
	//     "role":rolename,
	//     "permissions":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/role/permissions/grant", rbacAPI.GrantPermissionsToRole)

	// remove permission from specified role
	// Json params:
	// {
	//     "system":system,
	//     "role":rolename,
	//     "permission":permission
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/role/permissions/remove", rbacAPI.RemovePermissionFromRole)

	// register user
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "roles":[
	//         "roles1",
	//         "roles2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Post("/user", rbacAPI.RegisterUser)

	// unregister user
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Delete("/user", rbacAPI.UnregisterUser)

	// register user
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "new_roles":[
	//         "roles1",
	//         "roles2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user", rbacAPI.UpdateUser)

	// get user permission info
	// URL params: system, uid
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "user":{
	//         "system":system,
	//         "uid":uid,
	//         "roles":[
	//             "role1",
	//             "role2"
	//         ],
	//         "blacklist":[
	//             "permission1",
	//             "permission2"
	//         ],
	//         "whitelist":[
	//             "permission1",
	//             "permission2"
	//         ]
	//     }
	// }
	app.Get("/user", rbacAPI.GetUser)

	// get all roles of specified user by uid
	// URL params: system, uid
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "roles":[
	//         "role1",
	//         "role2"
	//     ]
	// }
	app.Get("/user/roles", rbacAPI.GetAllRolesByUID)

	// update roles of specified user by uid
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "roles":[
	//         "roles1",
	//         "roles2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/roles", rbacAPI.UpdateRoles)

	// add specified roles to user
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "roles":[
	//         "roles1",
	//         "roles2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/roles/add", rbacAPI.AddRoles)

	// remove specified role to user
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "role":role
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/roles/remove", rbacAPI.RemoveRoles)

	// get blacklist of specified user
	// URL params: system, uid
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "blacklist":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	app.Get("/user/blacklist", rbacAPI.GetBlackList)

	// add specified permission into blacklist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "permissions":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/blacklist/add", rbacAPI.AddToBlackList)

	// remove specified permission from blacklist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "permission":permission
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/blacklist/remove", rbacAPI.RemoveFromBlackList)

	// clear blacklist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/blacklist/clear", rbacAPI.ClearBlackList)

	// get whitelist of specified user
	// URL params: system, uid
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message,
	//     "whitelist":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	app.Get("/user/whitelist", rbacAPI.GetWhiteList)

	// update whitelist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "whitelist":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/whitelist", rbacAPI.UpdateWhiteList)

	// add specified permission into whitelist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "permissions":[
	//         "permission1",
	//         "permission2"
	//     ]
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/whitelist/add", rbacAPI.AddToWhiteList)

	// remove specified permission from whitelist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid,
	//     "permission":permission
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/whitelist/remove", rbacAPI.RemoveFromWhiteList)

	// clear whitelist
	// Json params:
	// {
	//     "system":system,
	//     "uid":uid
	// }
	//
	// Response
	// {
	//     "code": 0, // 0-success
	//     "message":message
	// }
	app.Put("/user/whitelist/clear", rbacAPI.ClearWhiteList)

	return nil
}

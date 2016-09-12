## RBAC Service 接口文档

[TOC]

> Note: 请求中用到的 system 表示业务方名称，uid 是业务方 uid，permission 表示权限名称，role 表示角色名称，blacklist 表示权限黑名单，whitelist 表示权限白名单

### 校验是否有指定权限

#### 权限校验流程

> 1. 查询权限`黑名单`，如果命中，则表示`无`相应权限，否则继续以下操作
> 2. 查询权限`白名单`，如果命中，则表示`有`相应权限，否则继续以下操作
> 3. 查询用户角色列表，并根据角色列表查询得到用户拥有的所有权限，如果包含指定权限，则校验通过，否则，检验失败

#### 请求

```
Get /authenticate?system={system}&&uid={uid}&&permission={permission}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "permit":true // true or false
}
```



### 注册权限

#### 请求

```
Post /permission

{
    "system":system,
    "name":name,
    "desc":description {option}
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 注销权限

#### 请求

```
Delete /permission

{
    "system":system,
    "name":name
    "desc":description {option}
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 获取所有已注册的权限列表

#### 请求

```
Get /permission?system={system}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "permissions":[
        {
            "system":system,
            "name":name
            "desc":desc
        },
        {
            "system":system,
            "name":name
            "desc":desc
        }
    ]
}
```

### 更改权限名称

#### 请求

```
Put /permission

{
    "system":system,
    "oldname":oldname,
    "newname":newname
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 注册角色

#### 请求

```
Post /role

{
    "system":system,
    "name":name,
    "desc":description {option}
    "permissions":[
        "permission1",
        "permission2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 注销角色

#### 请求

```
Delete /role

{
    "system":system,
    "name":name
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 查询角色信息

#### 请求

```
Get /role?system={system}&role={rolename}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "role":{
        "system":system,
        "name":name,
        "desc":desc,
        "permissions":[
            "permission1",
            "permission2"
        ]
    }
}
```

### 查询所有已注册角色

#### 请求

```
Get /role/all?system={system}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "roles":[
        {
            "system":system,
            "name":name,
            "desc":desc,
            "permissions":[
                "permission1",
                "permission2"
            ]
        },
        {
            "system":system,
            "name":name,
            "desc":desc,
            "permissions":[
                "permission1",
                "permission2"
            ]
        }
    ]
}
```

### 更改角色名称

#### 请求

```
Put /role

{
    "system":system,
    "oldname":oldname,
    "newname":newname,
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 查询指定角色包含的权限列表

#### 请求

```
Get /role/permissions?system={system}&role={rolename}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "permissions":[
        "permission1",
        "permission2"
    ]
}
```

### 给角色赋予权限

#### 请求

```
Put /role/permissions/grant

{
    "system":system,
    "role":rolename,
    "permissions":[
        "permission1",
        "permission2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 从角色中移除权限

#### 请求

```
Put /role/permissions/remove

{
    "system":system,
    "role":rolename,
    "permission":permission
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 绑定用户

#### 请求

```
Post /user

{
    "system":system,
    "uid":uid,
    "roles":[
        "roles1",
        "roles2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 解绑用户

#### 请求

```
Delete /user

{
    "system":system,
    "uid":uid
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 更新用户权限信息

#### 请求

```
Put /user

{
    "system":system,
    "uid":uid,
    "new_roles":[
        "roles1",
        "roles2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 查询用户信息

#### 请求

```
Get /user?system={system}&uid={uid}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "user":{
        "system":system,
        "uid":uid,
        "roles":[
            "role1",
            "role2"
        ],
        "blacklist":[
            "permission1",
            "permission2"
        ],
        "whitelist":[
            "permission1",
            "permission2"
        ]
    }
}
```

### 查询用户拥有的角色

#### 请求

```
Get /user/roles?system={system}&uid={uid}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "roles":[
        "role1",
        "role2"
    ]
}
```

### 更改用户角色信息

#### 请求

```
Put /user/roles

{
    "system":system,
    "uid":uid,
    "roles":[
        "roles1",
        "roles2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 给用户赋予角色

#### 请求

```
Put /user/roles/add

{
    "system":system,
    "uid":uid,
    "roles":[
        "roles1",
        "roles2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 移除用户的指定角色

#### 请求

```
Put /user/roles/remove

{
    "system":system,
    "uid":uid,
    "role":role
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 查询权限黑名单列表

#### 请求

```
Get /user/blacklist?system={system}&uid={uid}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "blacklist":[
        "permission1",
        "permission2"
    ]
}
```

### 添加指定权限到权限黑名单

#### 请求

```
Put /user/blacklist/add

{
    "system":system,
    "uid":uid,
    "permissions":[
        "permission1",
        "permission2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 将指定权限从权限黑名单中移除

#### 请求

```
Put /user/blacklist/remove

{
    "system":system,
    "uid":uid,
    "permission":permission
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 清空权限黑名单

#### 请求

```
Put /user/blacklist/clear

{
    "system":system,
    "uid":uid
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 查询权限白名单

#### 请求

```
Get /user/whitelist?system={system}&uid={uid}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message,
    "whitelist":[
        "permission1",
        "permission2"
    ]
}
```

### 更新权限白名单

#### 请求

```
Put /user/whitelist

{
    "system":system,
    "uid":uid,
    "whitelist":[
        "permission1",
        "permission2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 添加指定权限到权限白名单

#### 请求

```
Put /user/whitelist/add

{
    "system":system,
    "uid":uid,
    "permissions":[
        "permission1",
        "permission2"
    ]
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 将指定权限从权限白名单中移除

#### 请求

```
Put /user/whitelist/remove

{
    "system":system,
    "uid":uid,
    "permission":permission
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```

### 清空权限白名单

#### 请求

```
Put /user/whitelist/clear

{
    "system":system,
    "uid":uid
}
```

#### 响应

```
{
    "code": 0, // 0-success
    "message":message
}
```


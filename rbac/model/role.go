package model

// Role contains multi-permission
type Role struct {
	System      string   `json:"system" bson:"system"`
	Name        string   `json:"name" bson:"name"`
	Desc        string   `json:"desc" bson:"desc"`
	Permissions []string `json:"permissions" bson:"permissions"`
}

func NewRole(system, name, desc string, permissions ...string) *Role {
	return &Role{
		System:      system,
		Name:        name,
		Desc:        desc,
		Permissions: permissions,
	}
}

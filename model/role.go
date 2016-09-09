package model

// Role contains multi-permission
type Role struct {
	System      string   `json:"system" bson:"system" validate:"required"`
	Name        string   `json:"name" bson:"name" validate:"required"`
	Desc        string   `json:"desc" bson:"desc"`
	Permissions []string `json:"permissions" bson:"permissions" validate:"required"`
}

func NewRole(system, name, desc string, permissions ...string) *Role {
	return &Role{
		System:      system,
		Name:        name,
		Desc:        desc,
		Permissions: permissions,
	}
}

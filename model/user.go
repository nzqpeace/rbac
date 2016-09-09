package model

type UserPermModel struct {
	System    string   `json:"system" bson:"system" validate:"required"`
	UID       string   `json:"uid" bson:"uid" validate:"required"`
	Roles     []string `json:"roles" bson:"roles" validate:"required"`
	BlackList []string `json:"blacklist" bson:"blacklist"`
	WhiteList []string `json:"whitelist" bson:"whitelist"`
}

func NewUserPermModel(system, uid string, roles ...string) *UserPermModel {
	return &UserPermModel{
		System:    system,
		UID:       uid,
		Roles:     roles,
		BlackList: []string{},
		WhiteList: []string{},
	}
}

package model

type UserPermModel struct {
	System    string   `json:"-" bson:"system"`
	UID       string   `json:"-" bson:"uid"`
	Roles     []string `json:"roles" bson:"roles"`
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

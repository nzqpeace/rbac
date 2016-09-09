package model

// Permission exports `Name` `Tag` `Desc`
type Permission struct {
	System string `json:"system" bson:"system" validate:"required"`
	Name   string `json:"name" bson:"name" validate:"required"`
	Desc   string `json:"desc" bson:"desc"`
}

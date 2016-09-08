package model

// Permission exports `Name` `Tag` `Desc`
type Permission struct {
	System string `json:"system" bson:"system"`
	Name   string `json:"name" bson:"name"`
	Desc   string `json:"desc" bson:"desc"`
}

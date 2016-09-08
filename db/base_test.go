package db

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

var (
	conf       *MgoConf
	db         *DataBase
	collection *Base
)

type Cowshed struct {
	Name    string `json:"name" bson:"name"`
	Age     int    `json:"age" bson:"age"`
	Address string `json:"address" bson:"address"`
}

func init() {
	conf = &MgoConf{
		Url: "localhost/test",
	}

	var err error
	db, err = Init(conf)
	if err != nil {
		fmt.Println(err)
	}

	collection = NewBase(db, "Cowshed")
}

func TestMongo(t *testing.T) {
	m := &Cowshed{
		Name:    "zhiqiang",
		Age:     99,
		Address: "Pudong, Shanghai, China",
	}

	assert.Nil(t, collection.Insert(m))

	var p Cowshed
	// query
	assert.Nil(t, collection.Find(bson.M{"age": 99}, &p))
	assert.Equal(t, m.Name, p.Name)
	assert.Equal(t, m.Age, p.Age)
	assert.Equal(t, m.Address, p.Address)

	// check count
	count, err := collection.Count(bson.M{"age": 99})
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	// update document
	assert.Nil(t, collection.Update(bson.M{"age": 99}, bson.M{
		"$set": bson.M{
			"name": "qiniu",
		},
	}))

	// check value
	assert.Nil(t, collection.Find(bson.M{"age": 99}, &p))
	assert.NotEqual(t, m.Name, p.Name)
	assert.Equal(t, m.Age, p.Age)
	assert.Equal(t, m.Address, p.Address)

	// remove document
	assert.Nil(t, collection.Remove(bson.M{"age": 99}))

	// check count
	count, err = collection.Count(bson.M{"age": 99})
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

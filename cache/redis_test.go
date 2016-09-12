package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	config *RedisConfig
	r      *Redis
)

func init() {
	config = &RedisConfig{
		Address:       ":6379",
		Password:      "",
		DB:            1,
		MaxConn:       100,
		IdleTimeout:   60,
		RetryInterval: 3,
		RetryTimes:    0,
	}

	r = NewRedis(config)
}

func TestSet(t *testing.T) {
	key := "Project"
	err := r.SAdd(key, "Cowshed0", "Cowshed1")
	assert.Nil(t, err)

	exist, err := r.SIsMembers(key, "Cowshed0")
	assert.Nil(t, err)
	assert.True(t, exist)

	// check members
	members, err := r.SMembers(key)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(members))

	// check not exist element
	exist, err = r.SIsMembers(key, "Cowshed1")
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = r.SIsMembers(key, "NotExist")
	assert.Nil(t, err)
	assert.False(t, exist)

	// remove "Cowshed0", "Cowshed1"
	err = r.SRem(key, "Cowshed0", "Cowshed1")
	assert.Nil(t, err)

	exist, err = r.SIsMembers(key, "Cowshed0")
	assert.Nil(t, err)
	assert.False(t, exist)
}

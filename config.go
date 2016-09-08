package rbac

import (
	"github.com/nzqpeace/rbac/cache"
	"github.com/nzqpeace/rbac/db"
)

type RBACConfig struct {
	Redis *cache.RedisConfig
	Mgo   *db.MgoConf
}

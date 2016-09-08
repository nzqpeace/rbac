package rbac

import (
	"github.com/nzqpeace/rbac/rbac/cache"
	"github.com/nzqpeace/rbac/rbac/db"
)

type RBACConfig struct {
	Redis *cache.RedisConfig
	Mgo   *db.MgoConf
}

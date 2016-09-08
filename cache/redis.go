package cache

import (
	"errors"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

//Redis object
type Redis struct {
	pool *redis.Pool
}

//NewRedis initiates a new Redis instance
func NewRedis(config *RedisConfig) *Redis {
	return &Redis{
		pool: NewRedisPool(config),
	}
}

// Do wrapper of redis.Do
func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// NewRedisPool create a instance of redis pool
func NewRedisPool(config *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxActive:   config.MaxConn,
		MaxIdle:     config.MaxConn,
		Wait:        true,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.Address)
			if err != nil {
				return nil, err
			}

			masterAddr := config.Address
			skipParseInfoSentinel := false

			// check whether a sentinel
			sentinelInfo, err := redis.String(conn.Do("INFO", "SENTINEL"))
			if err != nil {
				if strings.Contains(err.Error(), "NOAUTH Authentication required") {
					skipParseInfoSentinel = true // this is a master redis, not sentinel
				} else {
					conn.Close()
					return nil, err
				}
			}

			// parse `info sentinel` for master redis's address
			if len(sentinelInfo) > 0 && !skipParseInfoSentinel {
				if index := strings.Index(sentinelInfo, "status=ok"); index != -1 {
					start := index + 18
					address := sentinelInfo[start:]
					end := strings.Index(address, ",")
					if end == -1 {
						conn.Close()
						return nil, errors.New("parse master redis address failed")
					}
					masterAddr = address[:end]
				} else {
					conn.Close()
					return nil, errors.New("couldn't found alive master redis")
				}
			}

			if masterAddr != config.Address {
				conn.Close() // close connection to sentinel

				// connect to master redis node
				conn, err = redis.Dial("tcp", masterAddr)
				if err != nil {
					return nil, err
				}
			}

			if _, err := conn.Do("AUTH", config.Password); err != nil {
				if strings.Contains(err.Error(), "invalid password") {
					conn.Close()
					return nil, err
				}
			}

			if _, err := conn.Do("SELECT", config.DB); err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
	}
}

// SAdd add members into redis set
func (r *Redis) SAdd(key string, members ...string) (err error) {
	if len(members) == 0 {
		return
	}

	var params []interface{}
	params = append(params, key)
	for _, m := range members {
		params = append(params, m)
	}

	_, err = r.Do("sadd", params...)
	return
}

// SRem remove members from redis set
func (r *Redis) SRem(key string, members ...string) (err error) {
	if len(members) == 0 {
		return
	}

	var params []interface{}
	params = append(params, key)
	for _, m := range members {
		params = append(params, m)
	}
	_, err = r.Do("srem", params...)
	return
}

// SMembers list all members at specified redis set
func (r *Redis) SMembers(key string) ([]string, error) {
	return redis.Strings(r.Do("smembers", key))
}

// SIsMembers check whether a member of specified set
func (r *Redis) SIsMembers(key, value string) (bool, error) {
	return redis.Bool(r.Do("sismember", key, value))
}

func (r *Redis) Exists(key string) (bool, error) {
	return redis.Bool(r.Do("exists", key))
}

func (r *Redis) FlushDB() {
	r.Do("flushdb")
}

// Del delete specified key from redis
func (r *Redis) Del(keys ...string) (bool, error) {
	if len(keys) == 0 {
		return true, nil
	}

	var params []interface{}
	for _, k := range keys {
		params = append(params, k)
	}
	return redis.Bool(r.Do("del", params...))
}

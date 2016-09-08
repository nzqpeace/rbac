package cache

//RedisConfig is redis's configuration
type RedisConfig struct {
	Address       string `json:"address"`
	Password      string `json:"password"`
	DB            int    `json:"db"`
	MaxConn       int    `json:"max_conn"`
	IdleTimeout   int    `json:"idle_timeout"`
	RetryInterval int    `json:"retry_interval"`
	RetryTimes    int    `json:"retry_times"`
}

func DefaultConfig() *RedisConfig {
	return &RedisConfig{
		Address:       "localhost:6379",
		Password:      "",
		DB:            0,
		MaxConn:       100,
		IdleTimeout:   60,
		RetryInterval: 0,
		RetryTimes:    0,
	}
}

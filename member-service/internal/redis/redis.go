package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	mlog "github.com/mike504110403/goutils/log"
)

var (
	RedisPool *redis.Pool
)

type Config struct {
	RedisServer   string
	RedisPassword string
	RedisDB       int
}

func Init(cfg Config) {
	// 初始化 Redis 连接池
	RedisPool = newRedisPool(cfg.RedisServer, cfg.RedisPassword, cfg.RedisDB)
	conn := RedisPool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		mlog.Fatal(fmt.Sprintf("Redis 初始化失敗: %v", err))
	}
}

func newRedisPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if db != 0 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			// 每次取用連線前檢查是否連線還在
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetRedisConn() redis.Conn {
	return RedisPool.Get()
}

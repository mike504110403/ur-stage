package gameoperatorhandler

import (
	"fmt"
	"time"

	"member_service/internal/locals"
	redisDriver "member_service/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	mlog "github.com/mike504110403/goutils/log"
)

func CheckGameOperator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser, err := locals.GetUserInfo(c)
		if err != nil {
			mlog.Error(fmt.Sprintf("Error getting user info: %v", err))
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"message": "未登入"})
		}

		isLocked, err := getRedisLock(localUser.MemberId, 30*time.Second)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"message": "内部服务器错误"})
		}
		if !isLocked {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(fiber.Map{"message": "操作过于频繁，请稍后再试"})
		}
		return c.Next()
	}
}

func getRedisLock(mid int, duration time.Duration) (bool, error) {
	key := fmt.Sprintf("member:lock:%d", mid)
	conn := redisDriver.GetRedisConn()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("SETNX", key, "locked"))
	if err != nil {
		return false, err
	}
	if reply == 1 {
		_, err = conn.Do("EXPIRE", key, int(duration.Seconds()))
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func ReleaseRedisLock(mid int) error {
	key := fmt.Sprintf("member:lock:%d", mid)
	conn := redisDriver.GetRedisConn()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

package cookies

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

var cfg Config

// Init : 設定cookies的config
func Init(c ...Config) {
	if len(c) != 0 {
		cfg = c[0]
	} else {
		cfg = defaultConfig
	}
}

// Get : 取得cookies
func (k Key) Get(c *fiber.Ctx) string {
	return c.Cookies(string(k))
}

// Set : 設定cookies
func (k Key) Set(c *fiber.Ctx, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     string(k),
		Value:    value,
		MaxAge:   cfg.MaxAge,
		HTTPOnly: cfg.HTTPOnly,
		Secure:   cfg.Secure,
		Domain:   cfg.Domain,
		SameSite: cfg.SameSite,
	})
}

// Clear : 清除cookies
func (k Key) Clear(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     string(k),
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: cfg.HTTPOnly,
		Secure:   cfg.Secure,
		Domain:   cfg.Domain,
		SameSite: cfg.SameSite,
	})
}

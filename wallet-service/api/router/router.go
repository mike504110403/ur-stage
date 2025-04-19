package router

import (
	"gitlab.com/gogogo2712128/common_moduals/fiber/handler/recoverhandler"
	"gitlab.com/gogogo2712128/common_moduals/fiber/handler/tracehandler"
	"gitlab.com/gogogo2712128/common_moduals/ilog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"

	v1 "wallet_service/api/handler/v1"
)

// Set : 設定全部的路由
func Set(r *fiber.App) error {
	r.Use(
		recoverhandler.New(),
		tracehandler.New(),
		func(c *fiber.Ctx) error { // 設定每次request建立一個log物件，並在最後處理或印出log
			reqUrl := c.Request().URI().RequestURI()
			logData := ilog.Basic().Trace(tracehandler.Get(c))
			defer func() {
				if c.Response().StatusCode() != fiber.StatusOK {
					if len(c.Request().Body()) > 200 {
						logData.Log(`[req][path->%v][body]%v`, c.Path(), string(c.Request().Body()[:200]))
					} else {
						logData.Log(`[req][path->%v][body]%v`, c.Path(), string(c.Request().Body()))
					}
					logData.Log(`[res->%v][path->%v] %v`, c.Response().StatusCode(), c.Path(), string(reqUrl))
				}
			}()
			return c.Next()
		},
		cors.New(
			cors.Config{
				AllowOrigins: "*",
			},
		),
		helmet.New(
			helmet.Config{
				XFrameOptions: "SAMEORIGIN",
			},
		),
	)

	v1.SetRouter(r)

	return nil
}

package auth

import (
	"database/sql"

	"member_service/internal/database"
	"member_service/internal/locals"
	codeValidator "member_service/internal/validator"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"gopkg.in/go-playground/validator.v9"

	jwthandler "member_service/api/router/middleware/jwthandler"
	"member_service/internal/access"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/auth")
	{
		g.Post("/login", loginHandler)                     // 登入
		g.Post("/register", regiterHandler)                // 註冊
		g.Post("/logout", jwthandler.New(), logoutHandler) // 登出
		g.Post("/validator", validatorHandler)             // 請求驗證碼
		g.Post("/reset", jwthandler.New(), ResetHandler)   // 重設密碼
	}
}

// 登入
func loginHandler(c *fiber.Ctx) error {
	req := LoginReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数错误"})
	}
	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求解析错误"})
	}

	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	l := LoginInit(LoginType(req.LoginType), db)
	// 取得會員資料
	err = l.GetMember(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"message": "账号或密码错误"})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	// 密碼驗證
	if valid := l.LoginValid(req.Password); !valid {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "账号或密码错误"})
	}

	return l.LoginSucess(c)
}

// 註冊
func regiterHandler(c *fiber.Ctx) error {
	req := RegisterReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数错误"})
	}
	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求解析错误"})
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	r := RegisterInit(RegisterType(req.RegisterType), db)

	// 檢查帳號是否存在
	exist, err := r.IsMemberExist(req.Username)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	if exist {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "账号已存在"})
	}
	// 註冊檢查
	err = r.RegisterValid(req.Username, req.Email, req.ValidateCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "注册验证失败"})
	}

	// 寫入註冊
	if err := r.RegisterRecord(
		req.Username,
		req.Password,
		req.Nickname,
		req.Name,
		req.Email,
		req.PromoCode,
	); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			// 根據錯誤碼識別唯一鍵約束違規
			if mysqlErr.Number == 1062 {
				return c.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{"message": "账号或email已存在"})
			}
		}
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return RegisterSucess(c)
}

// 登出
func logoutHandler(c *fiber.Ctx) error {
	access.Logout(c)
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "登出成功"})
}

// 請求驗證碼
func validatorHandler(c *fiber.Ctx) error {
	req := ValidatorReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数错误"})
	}

	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求解析错误"})
	}

	// 驗證碼發送
	if err := Validator(req); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "验证码已发送"})
}

// 重設密碼
func ResetHandler(c *fiber.Ctx) error {
	req := ResetReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数错误"})
	}

	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求解析错误"})
	}

	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	if validated := codeValidator.VerifyCode(localUser.UserName, req.ValidVal, req.Code); !validated {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "验证码错误"})
	}

	if err := ResetPassword(req.Password, localUser.MemberId); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "重设密码成功"})
}

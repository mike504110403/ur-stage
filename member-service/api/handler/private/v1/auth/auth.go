package auth

import (
	"database/sql"
	"member_service/instance/login"
	"member_service/instance/register"

	basicauth "gitlab.com/gogogo2712128/common_moduals/auth/basicAuth"

	"github.com/gofiber/fiber/v2"
)

// 登入
type Login interface {
	GetMember(string) error
	LoginValid(string) bool
	LoginSucess(*fiber.Ctx) error
}

type LoginType string

const (
	AccountLogin LoginType = "account"
)

// 註冊
type Register interface {
	IsMemberExist(string) (bool, error)
	RegisterValid(string, string, string) error
	RegisterRecord(string, string, string, string, string, *string) error
	RegisterSucess(*fiber.Ctx) error
}
type RegisterType string

const (
	AccountRegister RegisterType = "account"
)

var login_instance = &login.AccountLogin{}

func LoginInit(s LoginType, db *sql.DB) Login {
	if s == AccountLogin {
		basicauth.Init(basicauth.Config{
			UseEncodeType: basicauth.HashTypeSHA512,
			PasswordHash:  "QF2aJflsiiiFE25g",
		})
		login_instance.MemberDb = db
		return login_instance
	}
	return login_instance
}
func GetMember(u string) error {
	return login_instance.GetMember(u)
}
func LoginValid(p string) bool {
	return login_instance.LoginValid(p)
}
func LoginSucess(c *fiber.Ctx) error {
	return login_instance.LoginSucess(c)
}

var register_instance = &register.AccountRegister{}

func RegisterInit(s RegisterType, db *sql.DB) Register {
	if s == AccountRegister {
		basicauth.Init(basicauth.Config{
			UseEncodeType: basicauth.HashTypeSHA512,
			PasswordHash:  "QF2aJflsiiiFE25g",
		})

		register_instance.MemberDb = db
		return register_instance
	}
	basicauth.Init(basicauth.Config{
		UseEncodeType: basicauth.HashTypeSHA512,
		PasswordHash:  "QF2aJflsiiiFE25g",
	})
	return register_instance
}
func IsMemberExist(username string) (bool, error) {
	return register_instance.IsMemberExist(username)
}
func RegisterValid(username string, validVal string, code string) error {
	return register_instance.RegisterValid(username, validVal, code)
}
func RegisterRecord(username string, password string, nickname string, name string, email string, promo *string) error {
	return register_instance.RegisterRecord(username, password, nickname, name, email, promo)
}
func RegisterSucess(c *fiber.Ctx) error {
	return register_instance.RegisterSucess(c)
}

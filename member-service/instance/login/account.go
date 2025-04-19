package login

import (
	"database/sql"
	"member_service/internal/access"
	"member_service/internal/statuscode"

	basicauth "gitlab.com/gogogo2712128/common_moduals/auth/basicAuth"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
)

type AccountLogin struct {
	MemberDb *sql.DB
	Member   MemberLogin
}

type MemberLogin struct {
	MemberId    int    `db:"id"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	MemberState int    `db:"member_state"`
}

// 取得會員資料
func (a *AccountLogin) GetMember(username string) error {
	// 已啟用狀態
	memberEnable := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "enable",
	}
	memberAdmin := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "admin",
	}
	enableState, err := memberEnable.Get()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	adminState, err := memberAdmin.Get()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	queryStr := `
		SELECT id, username, password, member_state
		FROM Member
		WHERE username = ? AND ( member_state = ? OR member_state = ? )
	`
	err = a.MemberDb.QueryRow(queryStr, username, enableState, adminState).Scan(&a.Member.MemberId, &a.Member.Username, &a.Member.Password, &a.Member.MemberState)
	if err != nil {
		mlog.Error("資料庫查詢錯誤")
		return err
	}
	return nil
}

// 密碼驗證
func (a *AccountLogin) LoginValid(password string) bool {
	valid, err := basicauth.Check(password, a.Member.Password)
	if err != nil {
		mlog.Error(err.Error())
		return false
	}

	return valid
}

func (a *AccountLogin) LoginSucess(c *fiber.Ctx) error {
	// 登入成功後執行jwt登入成功流程
	token, err := access.LoginSuccess(c, a.Member.MemberId, a.Member.Username)
	if err != nil {
		return statuscode.JwtgenerateFail.ToRes().Err(err.Error()).ToErr()
	}
	c.Locals("message", "登录成功")
	return c.Status(fiber.StatusOK).
		JSON(token)
}

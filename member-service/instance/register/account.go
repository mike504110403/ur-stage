package register

import (
	"database/sql"
	"errors"
	"member_service/internal/access"
	"member_service/internal/database"
	"member_service/internal/statuscode"
	"member_service/internal/validator"

	"github.com/gofiber/fiber/v2"

	basicauth "gitlab.com/gogogo2712128/common_moduals/auth/basicAuth"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

type AccountRegister struct {
	MemberDb *sql.DB
	Member   struct {
		MemberId int    `db:"member_id"`
		Username string `db:"username"`
	}
}

// IsMemberExist 檢查會員是否存在 (檢查帳號是否重複)
func (r *AccountRegister) IsMemberExist(username string) (bool, error) {
	var exist int
	queryStr := `
		SELECT COUNT(*)
		FROM Member
		WHERE username = ? COLLATE utf8mb4_general_ci
	`
	err := r.MemberDb.QueryRow(queryStr, username).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist > 0, nil
}

// RegisterValid 註冊檢查 驗證
func (r *AccountRegister) RegisterValid(username string, validVal string, code string) error {
	if codeVerify := validator.VerifyCode(username, validVal, code); !codeVerify {
		return errors.New("驗證碼錯誤")
	}
	return nil
}

// RegisterRecord 寫入帳號資訊
func (r *AccountRegister) RegisterRecord(username string, password string, nickname string, name string, email string, promo *string) error {
	psHash, err := basicauth.HashPassword(password)
	if err != nil {
		return err
	}
	tx, err := r.MemberDb.Begin()
	if err != nil {
		return err
	}
	// 已啟用狀態
	memberEnable := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "enable",
	}
	state, err := memberEnable.Get()
	if err != nil {
		return err
	}
	// 寫入會員資訊
	insertStr := `
		INSERT INTO Member (username, password, member_state)
		VALUES (?, ?, ?)
	`
	mr, err := tx.Exec(insertStr, username, psHash, state)
	if err != nil {
		return err
	}
	// 紀錄會員資訊
	insertStr = `
		INSERT INTO MemberInfo (member_id, name, nick_name)
		VALUES (?, ?, ?)
	`
	memberId, err := mr.LastInsertId()
	if err != nil {
		return err
	}
	_, err = tx.Exec(insertStr, int(memberId), name, nickname)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 紀錄email
	insertStr = `
		INSERT INTO Security (member_id, email)
		VALUES (?, ?)
	`
	_, err = tx.Exec(insertStr, int(memberId), email)
	if err != nil {
		tx.Rollback()
		return err
	}
	if promo != nil {
		// 紀錄優惠碼
		insertStr = `
			INSERT INTO MemberPromo (member_id, promo_code)
			VALUES (?, ?)	
		`
		_, err = tx.Exec(insertStr, int(memberId), promo)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 建立錢包
	if err := prepareAccount(tx, int(memberId)); err != nil {
		tx.Rollback()
		return err
	}

	r.Member.MemberId = int(memberId)
	r.Member.Username = username
	return tx.Commit()
}

// prepareAccount 帳號準備
func prepareAccount(tx *sql.Tx, mid int) error {
	// 建立會員等級
	insertStr := `
		INSERT INTO MemberLevel (member_id, vip_level, member_level)
		VALUES (?, 0, 0)
		ON DUPLICATE KEY UPDATE member_id = member_id
	`

	_, err := tx.Exec(insertStr, mid)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertStr = `
		INSERT INTO Security (member_id)
		VALUES (?)
		ON DUPLICATE KEY UPDATE phone = phone, email = email
	`
	_, err = tx.Exec(insertStr, mid)
	if err != nil {
		tx.Rollback()
		return err
	}

	db, err := database.WALLET.DB()
	if err != nil {
		return err
	}
	// 建立錢包
	insertStr = `
		INSERT INTO Wallet (member_id, balance, lock_amount)
		VALUES (?, 0, 0)
		ON DUPLICATE KEY UPDATE member_id = member_id
	`
	_, err = db.Exec(insertStr, mid)
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

func (r *AccountRegister) RegisterSucess(c *fiber.Ctx) error {
	// 登入成功後執行jwt登入成功流程
	if _, err := access.LoginSuccess(c, r.Member.MemberId, r.Member.Username); err != nil {
		return statuscode.JwtgenerateFail.ToRes().Err(err.Error()).ToErr()
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "注册成功"})
}

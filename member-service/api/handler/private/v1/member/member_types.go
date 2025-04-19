package member

import "time"

// 會員狀態回應
type (
	MemberState string
	MemberRole  string
	StatusRes   struct {
		MemberId    int     `db:"member_id" json:"member_id"`       // 會員ID
		VipLevel    int     `db:"vip_level" json:"vip_level"`       // 會員VIP等級
		MemberLevel int     `db:"member_level" json:"member_level"` // 會員等級
		Balance     float64 `db:"balance" json:"balance"`           // 會員餘額
		LockBalance float64 `db:"lock_balance" json:"lock_balance"` // 會員凍結餘額
		Status      string  `json:"status" json:"status"`           // 會員狀態
		Role        string  `json:"role" json:"role"`               // 會員角色
	}
)

const (
	Enable  MemberState = "enable"
	Disable MemberState = "disable"
	Admin   MemberRole  = "admin"
	Normal  MemberRole  = "normal"
	Tester  MemberRole  = "tester"
)

// 會員基本資料回應
type (
	InfoRes struct {
		MemberId int       `json:"member_id"`  // 會員ID
		Basic    BasicInfo `json:"basic_info"` // 基本資料
		Security Security  `json:"security"`   // 帳號安全資料
	}
	BasicInfo struct {
		Username     string    `json:"username" db:"username"`
		Name         string    `json:"name" db:"name"`
		NickName     string    `json:"nick_name" db:"nick_name"`
		Gender       *string   `json:"gender" db:"gender"`
		Birthday     *string   `json:"birthday" db:"birthday"`
		RegisterDate time.Time `json:"register_date" db:"create_at"`
	}
	Security struct {
		Phone *string `json:"phone"`
		Email *string `json:"email"`
	}
)

// 會員交易紀錄請求
type (
	TransRecordReq struct {
		Type      *int   `json:"type"`
		StartDate string `json:"start_date" validate:"required"`
		EndDate   string `json:"end_date" validate:"required"`
	}
	TransRecordRes struct {
		TransType     string     `json:"trans_type"`
		TransDate     time.Time  `json:"trans_date"`
		OrderSeq      string     `json:"order_seq"`
		TransInfo     string     `json:"trans_info"`
		TransDoneDate *time.Time `json:"trans_done_date"`
		Amount        float64    `json:"amount"`
		Fee           float64    `json:"fee"`
		Status        string     `json:"status"`
	}
)

// 會員投注紀錄請求
type (
	BetRecordReq struct {
		GameType  *int   `json:"game_type"`
		StartDate string `json:"start_date" validate:"required"`
		EndDate   string `json:"end_date" validate:"required"`
	}
	BetRecordRes struct {
		GameType   string    `json:"game_type"`
		GameTypeId int       `json:"game_type_id" db:"game_type_id"`
		EffectBet  float64   `json:"effect_bet" db:"bet"`
		WinLose    float64   `json:"win_lose" db:"win_lose"`
		BetAt      time.Time `json:"bet_at" db:"bet_at"`
	}
)

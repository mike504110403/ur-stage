package auth

// 登入頁
type (
	// 登入請求
	LoginReq struct {
		Username  string  `json:"username" validate:"required"`
		Password  string  `json:"password" validate:"required"`
		LoginType string  `json:"login_type" validate:"required"`
		TgID      *string `json:"tg_id"`
	}
	// 登入回應
	LoginRes struct {
		Token string `json:"token"`
	}
)

// 註冊頁
type (
	// 註冊請求
	RegisterReq struct {
		Username     string  `json:"username" validate:"required"`
		Password     string  `json:"password" validate:"required"`
		Email        string  `json:"email" validate:"required"`
		ValidateCode string  `json:"validate_code" validate:"required"`
		Nickname     string  `json:"nickname" validate:"required"`
		Name         string  `json:"name" validate:"required"`
		PromoCode    *string `json:"promo_code"`
		RegisterType string  `json:"register_type" validate:"required"`
		TgID         *string `json:"tg_id"`
	}
)

type (
	// 驗證碼請求
	ValidatorReq struct {
		Username   string `json:"username" validate:"required"`
		ValidType  string `json:"valid_type" validate:"required"`
		ValidValue string `json:"valid_value" validate:"required"`
	}
)

type (
	ResetReq struct {
		Password string `json:"password" validate:"required"`    // 密碼
		ValidVal string `json:"valid_value" validate:"required"` // 驗證值
		Code     string `json:"code" validate:"required"`        // 驗證碼
	}
)

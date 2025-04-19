package validator

import (
	"sync"
	"time"
)

// 驗證碼結構
type VerificationCode struct {
	Code      string
	ExpiresAt time.Time
}

// 驗證碼暫存
var verificationCodes = struct {
	sync.RWMutex
	codes map[string]VerificationCode
}{
	codes: make(map[string]VerificationCode),
}

// 緩存驗證碼
func CacheVerificationCode(username string, validateVal string, code string) {
	verificationCodes.Lock()
	defer verificationCodes.Unlock()

	// 設置驗證碼過期時間（例如5分鐘後）
	expiration := time.Now().Add(5 * time.Minute)
	verificationCodes.codes[username+"-"+validateVal] = VerificationCode{
		Code:      code,
		ExpiresAt: expiration,
	}
}

// 驗證驗證碼
func VerifyCode(username string, validateVal string, inputCode string) bool {
	verificationCodes.RLock()
	defer verificationCodes.RUnlock()

	// 從cache中查找驗證碼
	if verification, ok := verificationCodes.codes[username+"-"+validateVal]; ok {
		if time.Now().Before(verification.ExpiresAt) && verification.Code == inputCode {
			return true
		}
	}
	return false
}

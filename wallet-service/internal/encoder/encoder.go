package encoder

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// Hmac384_Encoder  Hmac384 加密
func Hmac384_Encoder(paintext string, key []byte) (string, error) {
	// 要加密的消息
	message := []byte(paintext)
	hmacHash := hmac.New(sha512.New384, key)
	_, err := hmacHash.Write(message)
	if err != nil {
		return "", err
	}
	// 計算 HMAC 值
	hmacSum := hmacHash.Sum(nil)
	return hex.EncodeToString(hmacSum), nil
}

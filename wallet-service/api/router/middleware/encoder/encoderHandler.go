package encoderhandler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"

	mlog "github.com/mike504110403/goutils/log"

	"github.com/gofiber/fiber/v2"
)

// 預設的config
var cfg = Config{
	HeaderKeyCode: "Crypto-Code",
	HeaderKeySign: "Crypto-Sign",
	HeaderKeyTime: "Crypto-Time",
	ValidDuration: time.Second * 3,
	GetKeyMap: func() map[string]string {
		return make(map[string]string)
	},
	HookAfterDecode: func(c *fiber.Ctx, code string) error {
		return c.Next()
	},
}

type Config struct {
	HeaderKeyCode string
	HeaderKeySign string
	HeaderKeyTime string
	ValidDuration time.Duration
	// 比對Code和金鑰的清單
	GetKeyMap func() map[string]string
	// 解密完後的後續處理
	HookAfterDecode func(c *fiber.Ctx, code string) error
}

func Init(initCfg Config) {
	cfg = initCfg
}

// New : 加解密中間件
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := string(c.Request().Header.Peek(cfg.HeaderKeyCode))
		sign := string(c.Request().Header.Peek(cfg.HeaderKeySign))
		timestamp := string(c.Request().Header.Peek(cfg.HeaderKeyTime))
		key, hexKey := []byte{}, ""

		// 取得header中的指定名稱的值作為客戶的code，並比對出相對應的aes key
		if foundKey, isExist := cfg.GetKeyMap()[code]; !isExist {
			mlog.Error("Client Key 取得錯誤")
			return c.SendStatus(fiber.StatusBadRequest)
		} else {
			if keyByte, err := hex.DecodeString(foundKey); err != nil {
				mlog.Error("Client Key 解碼錯誤")
				return c.SendStatus(fiber.StatusBadRequest)
			} else {
				// 需要另外把hexkey裝載下來是因為驗簽是用hex編碼的key來當作金鑰的
				hexKey = foundKey
				key = keyByte
			}
		}

		// 進行驗簽
		if !ValidateSignature(sign, timestamp, code, []byte(hexKey)) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		// 將body進行解密
		if decryptedBody, err := Decrypt(c.Body(), key); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		} else {
			// 解開後將request body改成解開的內容
			c.Request().SetBody(decryptedBody)
			if err := cfg.HookAfterDecode(c, code); err != nil {
				return err
			}
		}

		// 回應內容透過一樣的方式加密
		defer func() {
			if c.Response().StatusCode() == fiber.StatusOK {
				if encryptedBody, err := Encrypt(c.Response().Body(), key); err != nil {
					c.Status(fiber.StatusBadRequest)
				} else {
					c.Response().SetBody(encryptedBody)
				}
			}
		}()

		return c.Next()
	}
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 填充以達到指定大小
	plaintext = pad(plaintext, aes.BlockSize)

	// 生成隨機IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 加密
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// 把IV串在前面
	fullCiphertext := append(iv, ciphertext...)

	// 轉Base64
	return []byte(base64.StdEncoding.EncodeToString(fullCiphertext)), nil
}

func Decrypt(encrypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Base64解密
	encryptedBytes, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		mlog.Error("Body 轉 base64 錯誤")
		return nil, err
	}

	// 依照BlockSize分割IV和密文
	iv := encryptedBytes[:aes.BlockSize]
	ciphertext := encryptedBytes[aes.BlockSize:]

	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 移除填充
	return unpad(plaintext)
}

// pad 使用PKCS#7填充
func pad(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	return append(src, padding...)
}

// unpad 删除PKCS#7填充
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("unpad error: input data is empty")
	}

	padLen := int(src[length-1])
	if padLen > length || padLen == 0 {
		return nil, fmt.Errorf("unpad error: invalid padding length")
	}

	// 檢查填充內容是否正確
	for _, v := range src[length-padLen:] {
		if int(v) != padLen {
			return nil, fmt.Errorf("unpad error: invalid padding content")
		}
	}

	return src[:length-padLen], nil
}

func CreateSignature(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidateSignature 驗證簽章
func ValidateSignature(receivedSignature string, receivedTimestamp string, receivedUUID string, secretKey []byte) bool {
	timestamp, err := strconv.ParseInt(receivedTimestamp, 10, 64)
	if err != nil {
		return false
	}
	receivedTime := time.Unix(timestamp, 0)
	// 先比對看看對方給的時間，已經過期就直接拒絕
	if time.Now().After(receivedTime.Add(cfg.ValidDuration)) || time.Now().Before(receivedTime) {
		mlog.Error("驗簽過期")
		return false
	}

	// 將uuid和timestamp拼起來做成簽章
	message := receivedUUID + receivedTimestamp
	expectedSignature := CreateSignature(message, secretKey)

	// 最後比對一下簽章是否正確
	return hmac.Equal([]byte(receivedSignature), []byte(expectedSignature))
}

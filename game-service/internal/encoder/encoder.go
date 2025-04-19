package encoder

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var cfg = Config{
	DesKey: "ecBXD5GM",
	DesIv:  "m8JW3pb9",
}

type Config struct {
	DesKey string
	DesIv  string
}

// Md5Encrypt : md5 加密
func Md5Encrypt(content string) string {
	md5Hash := md5.Sum([]byte(content))
	return hex.EncodeToString(md5Hash[:])
}

// 加密
func Encrypt(content []byte) (string, error) {
	block, err := des.NewCipher([]byte(cfg.DesKey))
	if err != nil {
		return "", err
	}

	// 填充以達到指定大小
	content = pad(content, des.BlockSize)

	// 加密
	ciphertext := make([]byte, len(content))
	mode := cipher.NewCBCEncrypter(block, []byte(cfg.DesIv))
	mode.CryptBlocks(ciphertext, content)

	// 把IV串在前面
	fullCiphertext := append([]byte(cfg.DesIv), ciphertext...)

	// 轉Base64
	return base64.StdEncoding.EncodeToString(fullCiphertext), nil
}

func Decrypt(encryptedContent []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Base64解密
	encryptedBytes, err := base64.StdEncoding.DecodeString(string(encryptedContent))
	if err != nil {
		return nil, err
	}

	// 依照BlockSize分割IV和密文
	iv := encryptedBytes[:des.BlockSize]
	ciphertext := encryptedBytes[des.BlockSize:]

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

package encoder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	mlog "github.com/mike504110403/goutils/log"
)

type CashEncryption struct {
	DesKey       string `json:"DesKey"`
	DesIV        string `json:"DesIV"`
	ClientSecret string `json:"ClientSecret"`
	ClientID     string `json:"ClientID"`
	SystemCode   string `json:"SystemCode"`
	Httpdomain   string `json:"Httpdomain"`
	WebID        string `json:"WebId"`
}

type HashConfig struct {
	HashKey    string
	HashIV     string
	Httpdomain string
	Prefix     string
}

type SecretConfig struct {
	AppID      string
	AppSecret  string
	Sign_key   string
	EBC_key    string
	Httpdomain string
}

type WgSportSecretInfo struct {
	Md5key     string
	Prefix     string
	Upusername string
	Httpdomain string
}

// mt_live加密
func (ce *CashEncryption) EncryptionData(content []byte) ([]byte, error) {
	block, err := des.NewCipher([]byte(ce.DesKey))
	if err != nil {
		return nil, err
	}
	// iv mode
	mode := cipher.NewCBCEncrypter(block, []byte(ce.DesIV))
	// 補字節
	content = pad(content, des.BlockSize)

	ciphertext := make([]byte, len(content))
	// 加密
	mode.CryptBlocks(ciphertext, content)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// mt_lottery加密
func (hc *HashConfig) EncryptionData(content []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(hc.HashKey))
	if err != nil {
		return nil, err
	}
	// iv mode
	mode := cipher.NewCBCEncrypter(block, []byte(hc.HashIV))
	// 補字節
	content = pad(content, aes.BlockSize)

	ciphertext := make([]byte, len(content))
	// 加密
	mode.CryptBlocks(ciphertext, content)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// 加簽
func (ce *CashEncryption) SignatureData(data []byte, timestamp int64) string {
	xdata := strconv.FormatInt(timestamp, 10) + ce.ClientSecret + ce.ClientID + string(data)
	hash := md5.Sum([]byte(xdata))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

// GB_Game MD5 加密
func SignatureData(data string) string {
	hash := md5.Sum([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	return strings.ToUpper(hashStr)
}

// GB_Game DES-ECB 加密
func DESEncryptECB(data []byte, key string) ([]byte, error) {
	block, err := des.NewCipher([]byte(key[:8]))
	if err != nil {
		return nil, err
	}

	data = pad(data, block.BlockSize())
	ciphertext := make([]byte, len(data))

	for bs, be := 0, block.BlockSize(); bs < len(data); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(ciphertext[bs:be], data[bs:be])
	}

	return ciphertext, nil
}

func DESEncryptECBBase64(data []byte, key string) (string, error) {
	ciphertext, err := DESEncryptECB(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// WG_Game MD5 加密
func SignData(params map[string]string) string {

	type kv struct {
		Key   string
		Value string
	}

	var orderedParams []kv
	for k, v := range params {
		if v != "" {
			orderedParams = append(orderedParams, kv{k, v})
		}
	}

	// 打印過濾後的鍵
	keys := make([]string, len(orderedParams))
	for i, kv := range orderedParams {
		keys[i] = kv.Key
	}
	fmt.Printf("Filtered keys: %v\n", keys)

	// 生成鍵值對的陣列
	keyValuePairs := make([]string, len(orderedParams))
	for i, kv := range orderedParams {
		keyValuePairs[i] = fmt.Sprintf("%s=%s", kv.Key, kv.Value)
	}
	queryString := strings.Join(keyValuePairs, "&")
	// 拼接 MD5 密钥
	queryString += "&md5key=" + "5N56LkgbNcwDjLaGEU9b"

	// 计算 MD5 签名
	hash := md5.Sum([]byte(queryString))
	return hex.EncodeToString(hash[:])
}

// WG_Sport MD5 加密
func WGSportSignData(params string) string {
	hash := md5.Sum([]byte(params))
	return hex.EncodeToString(hash[:])
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// DecryptData 解密資料
func (hc *HashConfig) DecryptData(encryptedData string) ([]byte, error) {

	// 解碼加密資料
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// 創建 AES 解密區塊
	block, err := aes.NewCipher([]byte(hc.HashKey))
	if err != nil {
		return nil, err
	}

	// 使用 CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, []byte(hc.HashIV))
	rawDataBytes := make([]byte, len(ciphertext))
	mode.CryptBlocks(rawDataBytes, ciphertext)

	// 去除補字節
	rawDataBytes = unpad(rawDataBytes)

	return rawDataBytes, nil
}

// unpad 去除補字節
func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// 組MD5字串
func GetMD5string(req interface{}) string {

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error("json.Marshal error")
		return ""
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		mlog.Error("json.Unmarshal error")
	}
	val := reflect.ValueOf(req)
	typeOfReq := val.Type()

	var orderedKeys []string
	for i := 0; i < val.NumField(); i++ {
		field := typeOfReq.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			orderedKeys = append(orderedKeys, jsonTag)
		}
	}

	// 按照結構字段的順序排序
	var orderedParams []KV
	for _, key := range orderedKeys {
		if value, exists := userMap[key]; exists && value != "" {
			orderedParams = append(orderedParams, KV{key, value})
		}
	}

	// 打印過濾後的鍵
	keys := make([]string, len(orderedParams))
	for i, kv := range orderedParams {
		keys[i] = kv.Key
	}

	// 生成鍵值對的陣列
	keyValuePairs := make([]string, len(orderedParams))
	for i, kv := range orderedParams {
		keyValuePairs[i] = fmt.Sprintf("%s=%s", kv.Key, kv.Value)
	}
	queryString := strings.Join(keyValuePairs, "&")

	return queryString
}

// DES-CBC 加密
func DESEncryptionDataCBC(content []byte, DesKey, DesIV string) (string, error) {
	block, err := des.NewCipher([]byte(DesKey))
	if err != nil {
		return "", err
	}

	padData := pad(content, block.BlockSize())
	ciphertext := make([]byte, len(padData))
	mode := cipher.NewCBCEncrypter(block, []byte(DesIV))
	mode.CryptBlocks(ciphertext, padData)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DES-CBC 解密
func DESDecryptCBC(cipherText []byte, key, iv string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return string(unpad(plaintext)), nil
}

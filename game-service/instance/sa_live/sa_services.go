package sa_live

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"game_service/internal/cachedata"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

var SA_LIVE_SECRECT_INFO = SaLiveSecrect_Info{}
var se = SA_LIVE_SECRECT_INFO.SE
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	saliveSecret, ok := cacheSecrect["sa_live"]
	if !ok {
		mlog.Error("sa_live secret not found")
		return
	}
	saliveId, ok := cacheId["sa_live"]
	if !ok {
		mlog.Error("sa_live id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(saliveId)
	if err != nil {
		mlog.Error(fmt.Sprintf("SA_Live取得代理商ID失敗: %s", err.Error()))
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	err = json.Unmarshal([]byte(saliveSecret), &SA_LIVE_SECRECT_INFO.SE)
	if err != nil {
		mlog.Error(fmt.Sprintf("SA_Live取得Secret_Info失敗: %s", err.Error()))
		return
	}

	SA_LIVE_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	se = SA_LIVE_SECRECT_INFO.SE
}

func apiCaller(q string, s string, ur string, handler func(r *fasthttp.Response) error) error {

	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(ur)

	// 設置請求正文
	req.Header.SetContentType("application/x-www-form-urlencoded")
	// 使用 url.Values 構建查詢參數
	data := url.Values{}
	data.Set("q", q)
	data.Set("s", s)
	req.SetBodyString(data.Encode())

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		return err
	}
	return handler(resp)
}

func ToQueryString(req interface{}) (string, error) {
	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	var kvPairs []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")

		// 獲取字段值
		value := val.Field(i).Interface()
		var valueStr string

		// 空或零值檢查
		if val.Field(i).IsZero() {
			continue
		}

		// 根據字段類型進行處理
		switch v := value.(type) {
		case string:
			valueStr = v
		case float64:
			if v == float64(int64(v)) {
				valueStr = fmt.Sprintf("%.0f", v) // 沒有小數位
			} else {
				valueStr = fmt.Sprintf("%.2f", v) // 顯示兩位小數
			}
		case int, int64:
			valueStr = fmt.Sprintf("%d", v) // 將 int 轉換為字符串
		default:
			return "", fmt.Errorf("unsupported field type: %T", v)
		}

		kvPairs = append(kvPairs, fmt.Sprintf("%s=%s", tag, valueStr))
	}
	// 排序
	sort.Strings(kvPairs)

	// 組字符串
	return strings.Join(kvPairs, "&"), nil
}

// 填充數據
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// DES 加密
func DESEncrypt(data, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	data = pad(data, block.BlockSize())
	ciphertext := make([]byte, len(data))

	mode := cipher.NewCBCEncrypter(block, key)
	mode.CryptBlocks(ciphertext, data)

	return ciphertext, nil
}

// 將加密數據轉換為 Base64 字串
func DESEncryptToBase64(data, key []byte) (string, error) {
	encrypted, err := DESEncrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// MD5 簽名
func BuildMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

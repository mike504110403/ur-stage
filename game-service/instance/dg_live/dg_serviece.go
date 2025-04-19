package dg_live

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"game_service/internal/cachedata"
	"strconv"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

var DG_LIVE_SECRECT_INFO = DGLiveSecrect_Info{}
var SE = DG_LIVE_SECRECT_INFO.SE
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	dgliveSecret, ok := cacheSecrect["dg_live"]
	if !ok {
		mlog.Error("dg_live secret not found")
		return
	}
	dgliveId, ok := cacheId["dg_live"]
	if !ok {
		mlog.Error("dg_live id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(dgliveId)
	if err != nil {
		mlog.Error(fmt.Sprintf("dg_live 取得代理商ID失敗: %s", err.Error()))
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	if err = json.Unmarshal([]byte(dgliveSecret), &DG_LIVE_SECRECT_INFO.SE); err != nil {
		mlog.Error(fmt.Sprintf("dg_live 取得 Secret_Info 失敗: %s", err.Error()))
		return
	}

	DG_LIVE_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	SE = DG_LIVE_SECRECT_INFO.SE
}

func apiCaller(uri API_URI, reqBody *[]byte, handler func(r *fasthttp.Response) error) error {
	if !DG_LIVE_SECRECT_INFO.RefreshTime.Before(time.Now()) || !DG_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(SE.URL + string(uri))

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := GenerateMD5(SE.API_ACCOUNT + SE.API_SECRET + timestamp)
	req.Header.Set("agent", SE.API_ACCOUNT)
	req.Header.Set("sign", sign)
	req.Header.Set("time", timestamp)

	// 設置請求正文
	req.Header.SetContentType("application/json")
	if reqBody != nil {
		req.SetBody(*reqBody)
	}

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
	if err := handler(resp); err != nil {
		mlog.Error(fmt.Sprintf("API 請求失敗: %s", err.Error()))
		return err
	}
	return nil
}

// 生成 MD5 散列值
func GenerateMD5(input string) string {
	// 將輸入的字符串轉換為字節數組
	data := []byte(input)

	// 計算 MD5 散列值
	hash := md5.Sum(data)

	// 將散列值轉換為十六進制字符串
	return hex.EncodeToString(hash[:])
}

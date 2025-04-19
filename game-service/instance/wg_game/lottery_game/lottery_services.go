package lottery_game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"game_service/internal/apicaller"
	encoder "game_service/pkg/encoder"

	"github.com/mike504110403/goutils/log"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/rand"
)

const Httpdomain = "https://apiwinneradm.dv2test.com"

func GetMember(db *sql.DB, mid int) (userData MemberGameAccount, err error) {

	// 檢查使用者是否存在
	queryStr := `
		SELECT M.username, mi.nick_name
		FROM Member AS M
		LEFT JOIN MemberInfo AS mi ON M.id = mi.member_id
		LEFT JOIN MemberGameAccount AS mga ON M.id = mga.member_id
		WHERE M.id = ?;
	`
	err = db.QueryRow(queryStr, mid).Scan(&userData.Username, &userData.Nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return userData, errors.New("查無資料")
		}
		return userData, err
	}

	return userData, nil
}

// 遊戲商會員建立
func createUserAPI(user CreateUserApiReq, url string) (CreateUserApiRes, error) {
	resp := CreateUserApiRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	// 将 JSON 字节数组转换为字符串
	jsonString := string(body)

	// 示例：按键值对拆分
	lines := strings.Split(jsonString, ",")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		mlog.Error(err.Error())
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func GameAccountExist(username string) (*CheckUserRes, error) {

	// 確認遊戲商帳號是否存在
	resp, err := checkUserAPI(
		CheckUserReq{
			ApiId:    "f300101",
			Username: "F3_" + username,
		},
		Httpdomain+string(LotteryCheckUserUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return nil, err
	}
	if resp.ErrorCode == "CLIENT_NOT_EXIST" {
		mlog.Error("遊戲商帳號不存在")
		return &resp, errors.New("遊戲商帳號不存在")
	}

	return &resp, nil
}

// 檢查遊戲商帳號是否存在
func checkUserAPI(user CheckUserReq, url string) (CheckUserRes, error) {
	resp := CheckUserRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 修改遊戲商帳號資訊
func editUserAPI(user EditUserReq, url string) (EditUserRes, error) {
	resp := EditUserRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 確認交易單號
func transferCheckAPI(user TransferCheckReq, url string) (TransferCheckRes, error) {
	resp := TransferCheckRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		mlog.Error(err.Error())
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 進入遊戲
func forwardGameAPI(user ForwardGameReq, url string) (ForwardGameRes, error) {
	resp := ForwardGameRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 踢出遊戲
func kickoutAPI(user KickUserReq, url string) (KickUserRes, error) {
	resp := KickUserRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 轉入點數
func transferInAPI(user TransferInReq, url string) (TransferInRes, error) {
	resp := TransferInRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 轉出點數
func transferOutAPI(user TransferOutReq, url string) (TransferOutRes, error) {
	resp := TransferOutRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 取得會員點數
func pointUserAPI(user PointUserReq, url string) (PointUserRes, error) {
	resp := PointUserRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 取得注單資訊
func buyListGetAPI(user BuyListGetApiReq, url string) (BuyListGetApiRes, error) {
	resp := BuyListGetApiRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 查詢彩票打賞注單
func getGiftListGetApiv2API(user GiftListGetApiv2Req, url string) (GiftListGetApiv2Res, error) {
	resp := GiftListGetApiv2Res{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 查詢代理線輸贏(彩票)
func proxyWinloseGetAPI(user ProxyWinloseGetReq, url string) (ProxyWinloseGetRes, error) {
	resp := ProxyWinloseGetRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 查詢會員輸贏
func userWinloseGetAPI(user UserWinloseGetReq, url string) (UserWinloseGetRes, error) {
	resp := UserWinloseGetRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 查詢注單開獎號
func buySingleGetAPI(user BuySingleGetApiReq, url string) (BuySingleGetApiRes, error) {
	resp := BuySingleGetApiRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// RegisterValid 註冊檢查 應檢查註冊資訊格式是否正確 並非檢查帳號是否存在
func AccountRegister(tx *sql.Tx, mid int, member MemberGameAccount) error {

	// 註冊會員
	err := createUser(tx, mid, MemberGameAccount{
		Password: createPassword(),
	})
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	resp, err := createUserAPI(
		CreateUserApiReq{
			ApiId:      "f300101",
			Username:   member.Username,
			Proxyname:  "f3_pr00101",
			Experience: "n",
		},
		Httpdomain+string(LotteryCreateUserUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		tx.Rollback()
		return err
	}
	if resp.ErrorCode != "OK" {
		mlog.Error("註冊失敗")
		tx.Rollback()
		return errors.New("註冊失敗")
	}

	tx.Commit()
	return nil
}

func createPassword() string {
	const passwordLength = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	password := make([]byte, passwordLength)
	for i := range password {
		password[i] = charset[seededRand.Intn(len(charset))]
	}
	fmt.Print("password: ", string(password))
	return string(password)
}

func createUser(tx *sql.Tx, mid int, newMember MemberGameAccount) error {
	// TODO:
	insertQuery := `
			INSERT INTO Member.MemberGameAccount (member_id, game_password)
			VALUES (?, ?)
			`
	_, err := tx.Exec(insertQuery, mid, newMember.Password)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// // 進入遊戲
func AccountLogin(login ForwardGameReq) (string, error) {
	resp, err := getURLTokenAPI(
		login,
		Httpdomain+string(LotteryForwardGameUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return "", err
	}
	if resp.ErrorCode == "OK" {
		return resp.URL, nil
	} else {
		return "", errors.New("進入遊戲失敗")
	}
}

// 取得URL Token
func getURLTokenAPI(user ForwardGameReq, url string) (ForwardGameRes, error) {
	resp := ForwardGameRes{}

	body, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		fmt.Println("error:", err)
	}

	// MD5加密
	signature := encoder.SignData(userMap)

	// 將簽名加入到 user 的 Sign 欄位中
	user.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err = json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			// if resp.Code != "00000" {
			// 	return errors.New(resp.Message)
			// }
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		log.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

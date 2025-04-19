package apollo

import (
	"database/sql"
	"errors"
	"game_service/internal/apicaller"
	"game_service/internal/database"

	Member "gitlab.com/gogogo2712128/common_moduals/dbModel/Member"

	"encoding/json"

	"fmt"
	"math/rand"
	"time"

	"github.com/mitchellh/mapstructure"

	mlog "github.com/mike504110403/goutils/log"

	"github.com/valyala/fasthttp"
)

// TODO: 以下的實作應該是要呼叫遊戲商API的地方
// 以及我們這邊相應的服務
// 不需要用方法來實作介面
// 直接用一般的function即可

var Httpdomain = "https://platform.apl-gaming.com/v1"

var headers = map[string]string{
	"Content-Type": "application/json",
	// 這邊應該是要放遊戲商的key
	"Authorization": "Bearer your_token_here",
}

func GetMember(db *sql.DB, username string) (Member.Member, error) {
	var member Member.Member
	var exist int

	// 檢查使用者是否存在
	queryStr := `
		SELECT COUNT(*)
		FROM MEMBER
		WHERE username = ?
	`
	err := db.QueryRow(queryStr, username).Scan(&exist)
	if err != nil {
		return member, err
	}
	if exist == 0 {
		return member, errors.New("user not found")
	}

	// 查詢使用者詳細資料
	queryStr = `
		SELECT username
		FROM MEMBER
		WHERE username = ?
	`
	row := db.QueryRow(queryStr, username)
	err = row.Scan(&member.Username, &member.Password, "")

	//err = row.Scan(&member.Username, &member.Password, &member.NickName)
	if err != nil {
		return member, err
	}

	return member, nil
}

// 檢查會員是否存在 (檢查帳號是否重複)
func GameAccountExist(username string) (*ExistRes, error) {
	// 確認遊戲商帳號是否存在
	resp, err := getAccountExistAPI(
		ExistReq{
			Username: username,
		},
		Httpdomain+string(AccountExistUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return nil, err
	}
	if resp.Result != "Y" {
		mlog.Error("遊戲商帳號不存在")
		return &resp, errors.New("遊戲商帳號不存在")
	}

	return &resp, nil
}

// 取得遊戲商帳號
func GetGameAccount(db *sql.DB, mid int) (MemberGameAccount, error) {
	var gameAccount MemberGameAccount

	queryStr := `
		SELECT username, password, nickname
		FROM MEMBER_GAME_ACCOUNT
		WHERE member_id = ?
	`
	row := db.QueryRow(queryStr, mid)

	err := row.Scan(&gameAccount.Username, &gameAccount.Password, &gameAccount.NickName)
	if err != nil {
		return gameAccount, err
	}

	return gameAccount, nil
}

// 確認遊戲商帳號是否存在
func getAccountExistAPI(existReq ExistReq, url string) (ExistRes, error) {
	resp := ExistRes{}
	// 發送POST請求並處理響應
	body, err := json.Marshal(existReq)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	// r 接回傳的response轉型
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		//TODO:轉struct
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Result != "Y" {
			mlog.Error("登入失敗")
			return fmt.Errorf("登入失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 檢查會員是否存在
// func getAccountExistByUserName(db *sql.DB, mid string) (string, error) {
// 	var username string
// 	queryStr := `
//         SELECT username
//         FROM MEMBER
//         WHERE member_id = ?
//     `
// 	err := db.QueryRow(queryStr, mid).Scan(&username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", errors.New("會員帳號不存在")
// 		}
// 		return "", err
// 	}
// 	return username, nil
// }

// RegisterValid 註冊檢查 應檢查註冊資訊格式是否正確 並非檢查帳號是否存在
func AccountRegister(tx *sql.Tx, member Member.Member) error {

	// 註冊會員
	err := registerUser(tx, MemberGameAccount{
		Username: member.Username,
		Password: createPassword(),
		// NickName: member.NickName,
		NickName: "",
	})
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	resp, err := getAccountRegisterAPI(
		member,
		Httpdomain+string(AccountRegisterUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		tx.Rollback()
		return err
	}
	if resp.Status != 1 {
		mlog.Error("註冊失敗")
		tx.Rollback()
		return errors.New("註冊失敗")
	}

	tx.Commit()
	return nil
}

func getAccountRegisterAPI(member Member.Member, url string) (RegisterRes, error) {
	resp := RegisterRes{}

	body, err := json.Marshal(member)
	if err != nil {
		mlog.Error(err.Error())
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("註冊失敗")
			return fmt.Errorf("註冊失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func registerUser(tx *sql.Tx, newMember MemberGameAccount) error {
	insertStr := `
	INSERT INTO Member_Game_Account (username, password, nickname)
	VALUES (?, ?, ?)
	`

	// 註冊會員
	_, err := tx.Exec(insertStr, newMember.Username, newMember.Password, newMember.NickName)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func createPassword() string {
	const passwordLength = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, passwordLength)
	for i := range password {
		password[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(password)
}

// 額度存提
func Deposit(input DepositInput) error {
	resp, err := getDepositAPI(
		DepositReq{
			Username: input.Username,
			Gametype: input.Gametype,
			Amt:      input.Amt,
			Refno:    input.Refno,
			Type:     input.Type,
		},
		Httpdomain+string(DepositUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != 1 {
		mlog.Error("存提失敗")
		return errors.New("存提失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	// 註冊會員
	err = getDepositUser(db, input)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getDepositAPI(req DepositReq, url string) (DepositRes, error) {
	resp := DepositRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("註冊失敗")
			return fmt.Errorf("註冊失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO DB
func getDepositUser(db *sql.DB, input DepositInput) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND gametype = ? AND amt = ? AND refno = ? AND type = ?
`
	err := db.QueryRow(queryStr, input.Username, input.Gametype, input.Amt, input.Refno, input.Type).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, gametype, amt, refno, type)
	VALUES (?, ?, ?, ?, ?)
`
	_, err = db.Exec(insertStr, input.Username, input.Gametype, input.Amt, input.Refno, input.Type)
	if err != nil {
		return err
	}

	return nil
}

// 取得額度
func GetQuota(username string, gametype string) error {
	resp, err := getQuotaAPI(
		GetquotaReq{
			Username: username,
			Gametype: gametype,
		},
		Httpdomain+string(GetQuotaUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != 1 {
		mlog.Error("存提失敗")
		return errors.New("存提失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	//
	err = getQuotaUser(db, username, gametype)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getQuotaAPI(req GetquotaReq, url string) (GetquotaRes, error) {
	resp := GetquotaRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("取得額度失敗")
			return fmt.Errorf("取得額度失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO DB
func getQuotaUser(db *sql.DB, username string, gametype string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND gametype = ?
`
	err := db.QueryRow(queryStr, username, gametype).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, gametype)
	VALUES (?, ?)
`
	_, err = db.Exec(insertStr, username, gametype)
	if err != nil {
		return err
	}

	return nil
}

// 進入遊戲
func AccountLogin(login Login) (string, error) {
	resp, err := getAccountLoginAPI(
		login,
		Httpdomain+string(AccountLoginUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return "", err
	}
	if resp.Status != 1 {
		return "", errors.New("進入遊戲失敗")
	}
	return resp.Result.Url, nil
}

func getAccountLoginAPI(login Login, url string) (LoginRes, error) {
	resp := LoginRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(login)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("進入遊戲失敗")
			return fmt.Errorf("進入遊戲失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func BetReportTransfer(data string) (map[string]interface{}, error) {
	var datas map[string]interface{}

	err := json.Unmarshal([]byte(data), &datas)
	if err != nil {
		mlog.Error(err.Error())
		return nil, err
	}

	return datas, nil
}

// 取得報表
func GetBetReport(betReportReq BetReportReq) error {
	resp, err := getBetReportAPI(
		BetReportReq{
			Agid:      1,
			Startdate: betReportReq.Startdate,
			Enddate:   betReportReq.Enddate,
			Gametype:  betReportReq.Gametype,
		},
		Httpdomain+string(GetBetReportUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != 1 {
		mlog.Error("取得報表失敗")
		return errors.New("取得報表失敗")
	}

	// 將 BetListResult 轉換為 JSON 字串
	betListJSON, err := json.Marshal(resp.Result.Lottory.BetList)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	data, err := BetReportTransfer(string(betListJSON))
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	// 遍歷 BetList，逐筆新增到資料庫
	for _, bet := range data {
		err := mapstructure.Decode(bet, &data)
		username := data["Username"].(string)
		betId := data["BetId"].(string)
		betTime := data["BetTime"].(time.Time)
		payOut := data["PayOut"].(float64)

		if err != nil {
			mlog.Error("無法轉換 data: " + err.Error())
			return errors.New("無法轉換 data")
		}
		err = getBetReport(
			db,
			username,
			betId,
			betTime,
			payOut,
			betListJSON,
		)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
	}
	return nil
}

func getBetReportAPI(req BetReportReq, url string) (BetReportRes, error) {
	resp := BetReportRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}

		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}

		if resp.Status != 1 {
			mlog.Error("取得報表失敗")
			return fmt.Errorf("取得報表失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func getBetReport(
	db *sql.DB,
	username string,
	betId string,
	betTime time.Time,
	payOut float64,
	data any) error {

	insertStr := `
		INSERT INTO WALLET 
		(
			username,
			betId,
			betTime,
			payOut,
			data
		)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		insertStr,
		username,
		betId,
		betTime,
		payOut,
		data,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetBetReportList(db *sql.DB, mid string, gametype string, starttime time.Time, endtime time.Time) ([]BetListResult, error) {
	queryStr := `
    SELECT *
    FROM GAME_RECORD
    WHERE mid = ? AND startdate >= ? AND enddate <= ? AND gametype = ?
`
	rows, err := db.Query(queryStr, mid, starttime, endtime, gametype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []BetListResult
	for rows.Next() {
		var report BetListResult
		if err := rows.Scan(&report); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

// 變更密碼
func AccountChangePassword(username string, password string) error {
	resp, err := getAccountChangePasswordAPI(
		AccountChangePasswordReq{
			Username: username,
			Password: password,
		},
		Httpdomain+string(AccountChangePasswordUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != 1 {
		mlog.Error("變更密碼失敗")
		return errors.New("變更密碼失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getAccountChangePasswordUser(db, username, password)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getAccountChangePasswordAPI(req AccountChangePasswordReq, url string) (AccountChangePasswordRes, error) {
	resp := AccountChangePasswordRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("變更密碼失敗")
			return fmt.Errorf("變更密碼失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO  getAccountChangePasswordDB
func getAccountChangePasswordUser(db *sql.DB, username string, password string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND password = ?
`
	err := db.QueryRow(queryStr, username, password).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, password)
	VALUES (?, ?)
`
	_, err = db.Exec(insertStr, username, password)
	if err != nil {
		return err
	}

	return nil
}

// 取得帳本
func GetAccountBook(username string, gametype string) error {
	resp, err := getAccountBookAPI(
		GetAccountBookReq{
			Username: username,
			Gametype: gametype,
		},
		Httpdomain+string(GetAccountBookUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != 1 {
		mlog.Error("取得帳本失敗")
		return errors.New("取得帳本失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getAccountBookUser(db, username, gametype)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getAccountBookAPI(req GetAccountBookReq, url string) (GetAccountBookRes, error) {
	resp := GetAccountBookRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("取得帳本失敗")
			return fmt.Errorf("取得帳本失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getAccountBookDB
func getAccountBookUser(db *sql.DB, username string, gametype string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND gametype = ?
`
	err := db.QueryRow(queryStr, username, gametype).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, gametype)
	VALUES (?, ?)
`
	_, err = db.Exec(insertStr, username, gametype)
	if err != nil {
		return err
	}

	return nil
}

// 取得期數結果
func GetPeriodResult(input PeriodResultInput) error {
	resp, err := getPeriodResultAPI(
		PeriodResultReq{
			Gametype: input.Gametype,
			Gamedate: input.Gamedate,
			Gamenum:  input.Gamenum,
		},
		Httpdomain+string(GetPeriodResultUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != 1 {
		mlog.Error("取得帳本失敗")
		return errors.New("取得帳本失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getPeriodResultUser(db, input)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getPeriodResultAPI(req PeriodResultReq, url string) (PeriodResultRes, error) {
	resp := PeriodResultRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != 1 {
			mlog.Error("取得帳本失敗")
			return fmt.Errorf("取得帳本失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func getPeriodResultUser(db *sql.DB, input PeriodResultInput) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE gametype = ? AND gamedate = ? AND gamenum = ?
`
	err := db.QueryRow(queryStr, input.Gametype, input.Gamedate, input.Gamenum).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (gametype, gamedate, gamenum)
	VALUES (?, ?, ?)
`
	_, err = db.Exec(insertStr, input.Gametype, input.Gamedate, input.Gamenum)
	if err != nil {
		return err
	}

	return nil
}

// 取得有效投注額
func GetRealReport(input RealReportInput) error {
	resp, err := getRealReportAPI(
		RealReportReq{
			Username:  input.Username,
			Startdate: input.Startdate,
			Enddate:   input.Enddate,
			Gametype:  input.Gametype,
		},
		Httpdomain+string(GetRealReportUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("取得有效投注失敗")
		return errors.New("取得有效投注失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getRealReportUser(db, input)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getRealReportAPI(req RealReportReq, url string) (RealReportRes, error) {
	resp := RealReportRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("取得有效投注失敗")
			return fmt.Errorf("取得有效投注失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func getRealReportUser(db *sql.DB, input RealReportInput) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND startdate = ? AND enddate = ? AND gametype = ?
`
	err := db.QueryRow(queryStr, input.Username, input.Startdate, input.Enddate, input.Gametype).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, startdate, enddate, gametype)
	VALUES (?, ?, ?, ?)
`
	_, err = db.Exec(insertStr, input.Username, input.Startdate, input.Enddate, input.Gametype)
	if err != nil {
		return err
	}

	return nil
}

// 交易編號檢查-是否存在
func CheckRefno(username string, refno string) error {
	resp, err := getCheckRefnoAPI(
		RefnoReq{
			Username: username,
			Refno:    refno,
		},
		Httpdomain+string(CheckRefnoUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("交易編號不存在")
		return errors.New("交易編號不存在")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getCheckRefnoDB(db, username, refno)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getCheckRefnoAPI(req RefnoReq, url string) (RefnoRes, error) {
	resp := RefnoRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("交易編號檢查失敗")
			return fmt.Errorf("交易編號檢查失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO: getCheckRefnoDB
func getCheckRefnoDB(db *sql.DB, username string, refno string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND refno = ?
`
	err := db.QueryRow(queryStr, username, refno).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, refno)
	VALUES (?, ?)
`
	_, err = db.Exec(insertStr, username, refno)
	if err != nil {
		return err
	}

	return nil
}

// 取總投注額資料
func GetTotalRealGold(input TotalRealGoldInput) error {
	resp, err := getTotalRealGoldAPI(
		TotalRealGoldReq{Gametype: input.Gametype,
			Startdate: input.Startdate,
			Enddate:   input.Enddate,
		},
		Httpdomain+string(GetTotalRealGoldUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("取得總投注額資料失敗")
		return errors.New("取得總投注額資料失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getTotalRealGoldDB(db, input)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getTotalRealGoldAPI(req TotalRealGoldReq, url string) (TotalRealGoldRes, error) {
	resp := TotalRealGoldRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("取得總投注額資料失敗")
			return fmt.Errorf("取得總投注額資料失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getRealGoldDB
func getTotalRealGoldDB(db *sql.DB, input TotalRealGoldInput) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE gametype = ? AND startdate = ? AND enddate = ?
`
	err := db.QueryRow(queryStr, input.Gametype, input.Startdate, input.Enddate).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (gametype, startdate, enddate)
	VALUES (?, ?, ?)	
`
	_, err = db.Exec(insertStr, input.Gametype, input.Startdate, input.Enddate)
	if err != nil {
		return err
	}

	return nil
}

// 取捕魚細單
func GetFhDetails(betid string) error {
	//TODO:可暫緩
	resp, err := getFhDetailsAPI(
		FhDetailsReq{
			Betid: betid,
		},
		Httpdomain+string(GetFhDetailsUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("取捕魚細單失敗")
		return errors.New("取捕魚細單失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getFhDetailsDB(db, betid)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getFhDetailsAPI(req FhDetailsReq, url string) (FhDetailsRes, error) {
	resp := FhDetailsRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("取捕魚細單失敗")
			return fmt.Errorf("取捕魚細單失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getFhDetailsDB
func getFhDetailsDB(db *sql.DB, betid string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE betid = ?
`
	err := db.QueryRow(queryStr, betid).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (betid)
	VALUES (?)
`
	_, err = db.Exec(insertStr, betid)
	if err != nil {
		return err
	}

	return nil
}

// 取得對帳報表
func GetRecReport(input RecReportInput) error {
	resp, err := getRecReportAPI(
		RecReportReq{Username: input.Username,
			Startdate: input.Startdate,
			Enddate:   input.Enddate,
		},
		Httpdomain+string(GetRecReportUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("取得對帳報表失敗")
		return errors.New("取得對帳報表失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getRecReportDB(db, input)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getRecReportAPI(req RecReportReq, url string) (RecReportRes, error) {
	resp := RecReportRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("取得對帳報表失敗")
			return fmt.Errorf("取得對帳報表失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getRecReportDB
func getRecReportDB(db *sql.DB, input RecReportInput) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND startdate = ? AND enddate = ?
`
	err := db.QueryRow(queryStr, input.Username, input.Startdate, input.Enddate).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, startdate, enddate)
	VALUES (?, ?, ?)
`
	_, err = db.Exec(insertStr, input.Username, input.Startdate, input.Enddate)
	if err != nil {
		return err
	}

	return nil
}

// 退出遊戲
func AccountLogout(username string) error {
	resp, err := getAccountLogoutAPI(
		AccountLogoutReq{
			Username: username,
		},
		Httpdomain+string(AccountLogoutUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("退出遊戲失敗")
		return errors.New("退出遊戲失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getAccountLogoutDB(db, username)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getAccountLogoutAPI(req AccountLogoutReq, url string) (AccountLogoutRes, error) {
	resp := AccountLogoutRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("退出遊戲失敗")
			return fmt.Errorf("退出遊戲失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getAccountLogoutDB
func getAccountLogoutDB(db *sql.DB, username string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ?
`
	err := db.QueryRow(queryStr, username).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username)
	VALUES (?)
`
	_, err = db.Exec(insertStr, username)
	if err != nil {
		return err
	}

	return nil
}

// 檢查在線狀態
func CheckOnlineStatus(username string, gametype string) error {
	//TODO:目前接口只支援捕魚FH，可暫緩
	resp, err := getCheckOnlineStatusAPI(
		OnlineStatusReq{Username: username,
			Gametype: gametype,
		},
		Httpdomain+string(CheckOnlineStatusUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("查詢失敗")
		return errors.New("查詢失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getCheckOnlineStatusDB(db, username, gametype)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getCheckOnlineStatusAPI(req OnlineStatusReq, url string) (OnlineStatusRes, error) {
	resp := OnlineStatusRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("檢查在線狀態失敗")
			return fmt.Errorf("檢查在線狀態失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getOnlineStatusDB
func getCheckOnlineStatusDB(db *sql.DB, username string, gametype string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND gametype = ?
`
	err := db.QueryRow(queryStr, username, gametype).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, gametype)
	VALUES (?, ?)
`
	_, err = db.Exec(insertStr, username, gametype)
	if err != nil {
		return err
	}

	return nil
}

// 檢查盤口狀態
func GetLineStatus(username string) error {
	resp, err := getLineStatusAPI(
		GetLineReq{
			Username: username,
		},
		Httpdomain+string(GetLineStatusUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("查詢失敗")
		return errors.New("查詢失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getLineStatusDB(db, username)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getLineStatusAPI(req GetLineReq, url string) (GetLineRes, error) {
	resp := GetLineRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("檢查盤口狀態失敗")
			return fmt.Errorf("檢查盤口狀態失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getLineStatusDB
func getLineStatusDB(db *sql.DB, username string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ?
`
	err := db.QueryRow(queryStr, username).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username)
	VALUES (?)
`
	_, err = db.Exec(insertStr, username)
	if err != nil {
		return err
	}

	return nil
}

// 預存提接口
func PreDeposit(input PreDepositInput) error {
	resp, err := getPreDepositAPI(
		PreDepositReq{
			Username: input.Username,
			Gametype: input.Gametype,
			Amt:      input.Amt,
			Refno:    input.Refno,
			Type:     input.Type,
		},
		Httpdomain+string(PreDepositUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("預存提接口失敗")
		return errors.New("預存提接口失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getPreDepositDB(db, input)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getPreDepositAPI(req PreDepositReq, url string) (PreDepositRes, error) {
	resp := PreDepositRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("預存提接口失敗")
			return fmt.Errorf("預存提接口失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getPreDepositDB
func getPreDepositDB(db *sql.DB, input PreDepositInput) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE username = ? AND gametype = ? AND amt = ? AND refno = ? AND type = ?
`
	err := db.QueryRow(queryStr, input.Username, input.Gametype, input.Amt, input.Refno, input.Type).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (username, gametype, amt, refno, type)
	VALUES (?, ?, ?, ?, ?)
`
	_, err = db.Exec(insertStr, input.Username, input.Gametype, input.Amt, input.Refno, input.Type)
	if err != nil {
		return err
	}

	return nil
}

// 確認預存提接口
func CheckDeposit(preno string) error {
	resp, err := getCheckDepositAPI(
		CheckDepositReq{
			Preno: preno,
		},
		Httpdomain+string(CheckDepositUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if resp.Status != "1" {
		mlog.Error("確認預存提接口失敗")
		return errors.New("確認預存提接口失敗")
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	err = getCheckDepositDB(db, preno)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

func getCheckDepositAPI(req CheckDepositReq, url string) (CheckDepositRes, error) {
	resp := CheckDepositRes{}

	// 發送POST請求並處理響應
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
	}

	// 定義回調函數
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(err.Error())
			return err
		}
		responseBody := r.Body()
		// 將 responseBody 轉換為字串進行判斷
		err := json.Unmarshal(responseBody, &resp)
		if err != nil {
			mlog.Error(err.Error())
			return err
		}
		if resp.Status != "1" {
			mlog.Error("確認預存提接口失敗")
			return fmt.Errorf("確認預存提接口失敗")
		}
		return nil
	}
	err = apicaller.SendPostRequest(url, headers, body, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// TODO getCheckDepositDB
func getCheckDepositDB(db *sql.DB, preno string) error {
	exist := 0
	queryStr := `
	SELECT COUNT(*)
	FROM WALLET
	WHERE preno = ?
`
	err := db.QueryRow(queryStr, preno).Scan(&exist)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("交易編號已存在")
	}

	insertStr := `
	INSERT INTO WALLET (preno)
	VALUES (?)
`
	_, err = db.Exec(insertStr, preno)
	if err != nil {
		return err
	}

	return nil
}

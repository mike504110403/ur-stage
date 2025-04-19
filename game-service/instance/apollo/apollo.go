package apollo

import (
	"game_service/internal/database"
	"time"

	mlog "github.com/mike504110403/goutils/log"
)

// TODO: 這邊才是對“遊戲”介面的實作 層次要抽開

// Apollo apollo遊戲商
type ApolloGamer struct {
}

// 進入遊戲
func (a *ApolloGamer) JoinGame(mid string, gametype string) (string, error) {
	res := ""

	db, err := database.MEMBER.DB()
	if err != nil {
		return res, err
	}

	member, err := GetMember(db, mid)
	if err != nil {
		return res, err
	}

	tx, err := database.MEMBER.TX()
	if err != nil {
		mlog.Error(err.Error())
		return res, err
	}
	// 確認帳號是否存在
	gameExistRes, err := GameAccountExist(member.Username)
	if err != nil {
		// 其他類型錯誤
		if gameExistRes == nil {
			return res, err
		}
		// 帳號不存在
		err := AccountRegister(tx, member)
		if err != nil {
			mlog.Error(err.Error())
			return res, err
		}
	}

	gameAccount, err := GetGameAccount(db, member.Id)
	if err != nil {
		mlog.Error(err.Error())
		return res, err
	}

	if url, err := AccountLogin(Login{
		Username: gameAccount.Username,
		Password: gameAccount.Password,
		Gametype: gametype,
	}); err != nil {
		return "", err
	} else {
		res = url
	}

	return res, nil

}

// 離開遊戲
func (a *ApolloGamer) LeaveGame() error {
	return nil
}

// 積分轉入
func (a *ApolloGamer) PointIn() error {
	return nil
}

// 積分轉出
func (a *ApolloGamer) PointOut() error {
	return nil
}

// 前端取得報表
func (a *ApolloGamer) BetRecord(mid string, gametype string, starttime time.Time, endtime time.Time) ([]BetListResult, error) {
	var res []BetListResult
	// 假設你有一個已經初始化的資料庫連接
	db, err := database.MEMBER.DB()
	if err != nil {
		return res, err
	}
	res, err = GetBetReportList(db, mid, gametype, starttime, endtime)
	if err != nil {
		return res, err
	}

	return res, nil
}

// 取得遊戲商BetRecord，並且記錄在DB
func (a *ApolloGamer) GetBetRecord(agid int, starttime time.Time, endtime time.Time) {

}

package wg_game

// TODO: 這邊才是對“遊戲”介面的實作 層次要抽開

// WG_Game遊戲商
type WGGamer struct {
}

// 進入遊戲
// func (m *WGGamer) JoinGame(mid string) (string, error) {
// 	res := ""

// 	db, err := database.MEMBER.DB()
// 	if err != nil {
// 		return res, err
// 	}

// 	// 將 mid 轉換成 int
// 	midInt, err := strconv.Atoi(mid)
// 	if err != nil {
// 		return res, err
// 	}

// 	member, err := GetMember(db, midInt)
// 	if err != nil {
// 		return res, err
// 	}

// 	tx, err := database.MEMBER.TX()
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return res, err
// 	}

// 	// 	// 確認遊戲商帳號是否存在
// 	// 	gameExistRes, err := GameAccountExist(member.UserName)
// 	// 	if err != nil || gameExistRes.Data.Result == 0 {
// 	// 		// 其他類型錯誤
// 	// 		if gameExistRes == nil {
// 	// 			return res, err
// 	// 		}
// 	// 		// 帳號不存在
// 	// 		err := AccountRegister(tx, midInt, member)
// 	// 		if err != nil {
// 	// 			mlog.Error(err.Error())
// 	// 			return res, err
// 	// 		}
// 	// 	}

// 	// 	gameAccount, err := GetGameAccount(db, midInt)
// 	// 	if err != nil {
// 	// 		mlog.Error(err.Error())
// 	// 		return res, err
// 	// 	}

// 	// 	if url, err := AccountLogin(GetURLTokenReq{
// 	// 		SystemCode: "THGKE9JKS9",
// 	// 		WebId:      "GKTEST01",
// 	// 		UserId:     gameAccount.UserName,
// 	// 		Language:   "zh-TW",
// 	// 	}); err != nil {
// 	// 		return "", err
// 	// 	} else {
// 	// 		res = url
// 	// 	}

// 	return res, nil

// }

// // // 離開遊戲
// // func (m *WGGamer) LeaveGame() error {
// // 	return nil
// // }

// // // 積分轉入
// // func (m *WGGamer) PointIn() error {
// // 	return nil
// // }

// // // 積分轉出
// // func (m *WGGamer) PointOut() error {
// // 	return nil
// // }

// // // 前端取得報表
// // // func (m *MTGamer) BetRecord(mid string, gametype string, starttime time.Time, endtime time.Time) ([]BetListResult, error) {
// // // 	var res []BetListResult
// // // 	// 假設你有一個已經初始化的資料庫連接
// // // 	db, err := database.MEMBER.DB()
// // // 	if err != nil {
// // // 		return res, err
// // // 	}
// // // 	res, err = GetBetReportList(db, mid, gametype, starttime, endtime)
// // // 	if err != nil {
// // // 		return res, err
// // // 	}

// // // 	return res, nil
// // // }

// // // 取得遊戲商BetRecord，並且記錄在DB
// // func (m *WGGamer) GetBetRecord(agid int, starttime time.Time, endtime time.Time) {

// // }

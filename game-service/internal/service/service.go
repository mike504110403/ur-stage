package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"game_service/internal/database"
	"strings"

	mlog "github.com/mike504110403/goutils/log"
)

// 取得會員遊戲帳號資訊
func GetMember(db *sql.DB, agentId int, mid int) (userData MemberGameAccountInfo, err error) {
	mlog.Info(fmt.Sprintf("開始查詢會員資料 MemberID: %d, AgentID: %d", mid, agentId))
	userData.MemberId = mid
	queryStr := `
		SELECT m.username, mi.nick_name
		FROM Member AS m
		JOIN MemberInfo AS mi ON m.id = mi.member_id
		WHERE m.id = ?;
	`
	err = db.QueryRow(queryStr, mid).Scan(&userData.UserName, &userData.NickName)
	if err != nil {
		mlog.Error(fmt.Sprintf("查詢會員基本資料失敗 MemberID: %d, Error: %s", mid, err.Error()))
		return userData, err
	}
	mlog.Info(fmt.Sprintf("成功查詢會員基本資料 MemberID: %d, Username: %s, NickName: %s",
		mid, userData.UserName, userData.NickName))
	// 檢查使用者是否存在
	queryStr = `
		SELECT game_password, game_agent_id
		FROM MemberGameAccount
		WHERE member_id = ? AND game_agent_id = ?;
	`
	err = db.QueryRow(queryStr, mid, agentId).Scan(&userData.GamePassword, &userData.GameAgentId)
	if err != nil {
		if err == sql.ErrNoRows {
			mlog.Info(fmt.Sprintf("會員遊戲帳號不存在 MemberID: %d, AgentID: %d", mid, agentId))
			return userData, nil
		}
		mlog.Error(fmt.Sprintf("查詢會員遊戲帳號失敗 MemberID: %d, AgentID: %d, Error: %s",
			mid, agentId, err.Error()))
		return userData, err
	}
	mlog.Info(fmt.Sprintf("成功查詢會員完整資料 MemberID: %d, Username: %s, GameAgentID: %d",
		mid, userData.UserName, userData.GameAgentId))
	return userData, nil
}

func PrepareMember(db *sql.DB) (*sql.Stmt, error) {
	// 檢查使用者是否存在
	queryStr := `
		SELECT M.username, mi.nick_name, COALESCE(mga.game_password, ''), COALESCE(mga.game_agent_id, 0)
		FROM Member AS M
		LEFT JOIN MemberInfo AS mi ON M.id = mi.member_id
		LEFT JOIN MemberGameAccount AS mga ON M.id = mga.member_id AND mga.game_agent_id = ?
		WHERE M.id = ?;
	`
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		mlog.Error(fmt.Sprintf("stmt異常: %s", err.Error()))
		return nil, err
	}
	return stmt, nil
}

// DB 建立會員 帳號
func CreateUser(tx *sql.Tx, mid int, newMember MemberGameAccountInfo) error {
	upsertStr := `
		INSERT INTO MemberGameAccount (member_id, game_agent_id, game_password)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE game_password = ?;
	`
	_, err := tx.Exec(upsertStr, mid, newMember.GameAgentId, newMember.GamePassword, newMember.GamePassword)
	if err != nil {
		mlog.Error(fmt.Sprintf("會員註冊失敗: %s", err.Error()))
		tx.Rollback()
		return err
	}

	return nil
}

func CreateGameUser(tx *sql.Tx, mid int, member MemberGameAccountInfo) error {
	upsertStr := `
		INSERT INTO MemberGameAccount (member_id, game_agent_id, game_password)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE game_password = ?;
	`
	_, err := tx.Exec(upsertStr, mid, member.GameAgentId, member.GamePassword, member.GamePassword)
	if err != nil {
		mlog.Error(fmt.Sprintf("註冊寫入失敗: %s", err.Error()))
		tx.Rollback()
		return err
	}
	mlog.Info(fmt.Sprintf("密碼更新成功: %s", member.GamePassword))

	return nil
}

// 抽成共用方法
func CreateDBUser(tx *sql.Tx, mid int, newMember MemberGameAccountInfo, agentId int) error {
	upsertStr := `
		INSERT INTO MemberGameAccount (member_id, game_agent_id, game_password)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE game_password = ?;
	`
	_, err := tx.Exec(upsertStr, mid, agentId, newMember.GamePassword, newMember.GamePassword)
	if err != nil {
		mlog.Error(fmt.Sprintf("GB_ELEC註冊寫入失敗: %s", err.Error()))
		tx.Rollback()
		return err
	}
	mlog.Info(fmt.Sprintf("密碼更新成功: %s", newMember.GamePassword))

	return nil
}

func WriteBetRecord(db *sql.DB, data string, agent int) error {
	// var DBbetRecord BetRecord

	dbs, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	upsertStr := `
	INSERT INTO Game.BetRecord (member_id, agent_id, bet_unique, game_type_id, bet_at, bet, win_lose, bet_info)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	slecStr := `
	SELECT id
	FROM Member
	WHERE username = ? AND member_state = ?;
	`
	var id int

	stmt, err := db.Prepare(upsertStr)
	if err != nil {
		mlog.Error(fmt.Sprintf("Prepare失敗: %s", err.Error()))
		return err
	}
	defer stmt.Close()

	switch agent {
	case 11:
		var RecordData BuyListGetRes
		err := json.Unmarshal([]byte(data), &RecordData)
		if err != nil {
			mlog.Error(fmt.Sprintf("WG_Sport注單解析失敗: %s", err.Error()))
			return err
		}
		for _, record := range RecordData.Data {
			// 使用 strings.Split 將 record.UserID 分割
			parts := strings.Split(record.Username, "_")
			username := ""
			if len(parts) > 1 {
				username = parts[1] // 取出 - 後面的部分
			} else {
				// 處理沒有 - 的情況
				username = record.Username
			}
			err = dbs.QueryRow(slecStr, username, 2).Scan(&id)
			if err != nil {
				mlog.Error(fmt.Sprintf("WG_Sport無此玩家: %s", err.Error()))
				return err
			}
			recordJson, err := json.Marshal(record)
			if err != nil {
				mlog.Error(fmt.Sprintf("WG_Sport_record解析失敗: %s", err.Error()))
				return err
			}
			_, err = stmt.Exec(id, agent, record.Id, 18, record.InsTime, record.Gold, record.Result, recordJson)
			if err != nil {
				mlog.Error(fmt.Sprintf("WG_Sport注單寫入失敗: %s", err.Error()))
				return err
			}
		}
	}

	return nil
}

// 寫入注單
func WriteBetRecord2(tx *sql.Tx, record []BetRecord, handler func() error) error {
	upsertStr := `
	INSERT INTO Game.BetRecord (
		member_id,
		agent_id,
		bet_unique,
		game_type_id,
		bet_at,
		bet,
		effect_bet,
		win_lose,
		bet_info
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(upsertStr)
	if err != nil {
		mlog.Error(fmt.Sprintf("Prepare失敗: %s", err.Error()))
		return err
	}
	defer stmt.Close()

	for _, r := range record {
		if _, err = stmt.Exec(
			r.MemberId,
			r.AgentId,
			r.BetUnique,
			r.GameTypeId,
			r.BetAt,
			r.Bet,
			r.EffectBet,
			r.WinLose,
			r.BetInfo,
		); err != nil {
			return err
		}
	}
	return handler()
}

// func dbBetRecord(stmt *sql.Stmt, record interface{}, agent int) error {
// 	_, err := stmt.Exec(1, agent, record.No, gameTypeId, betTime, record.BetTotal, record.BetValid, record)
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return err
// 	}
// }

func GetBalance(mid int) (float64, error) {
	db, err := database.WALLET.DB()
	if err != nil {
		mlog.Error(fmt.Sprintf("資料庫連線失敗: %s", err.Error()))
		return 0, err
	}
	var balance float64
	queryStr := `
		SELECT balance
		FROM Wallet
		WHERE member_id = ?;
	`
	err = db.QueryRow(queryStr, mid).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			mlog.Error(fmt.Sprintf("查無資料: %s", err.Error()))
			return 0, errors.New("查無資料")
		}
		mlog.Error(fmt.Sprintf("資料庫內部錯誤: %s", err.Error()))
		return 0, err
	}
	return balance, nil
}

// func InsertAmount(mid int) error {
// 	db, err := database.WALLET.DB()
// 	if err != nil {
// 		return err
// 	}
// 	insertStr := `
// 		INSERT INTO Wallet (member_id)
// 		VALUES (?, ?, 0, NOW(), NOW());
// 	`
// 	_, err = db.Exec(insertStr, mid)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

package transfer

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// 取得會員錢包資訊
func getMemberWallet(rdb redis.Conn, memberID int) (MemberWallet, error) {
	key := fmt.Sprintf("member_wallet_balance:%d", memberID)

	data, err := redis.Bytes(rdb.Do("GET", key))
	if err != nil {
		return MemberWallet{}, err
	}

	var wallet MemberWallet
	err = json.Unmarshal(data, &wallet)
	if err != nil {
		return MemberWallet{}, err
	}

	return wallet, nil
}

// 設定會員錢包資訊
func setMemberWallet(rdb redis.Conn, memberID int, wallet MemberWallet) error {
	key := fmt.Sprintf("member_wallet_balance:%d", memberID)

	jsonData, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	_, err = rdb.Do("SET", key, jsonData)
	if err != nil {
		return err
	}

	return nil
}

// 主錢包到小錢包轉帳
func transferMainToSub(rdb redis.Conn, memberId, gameId int, amount float64) error {
	// 會員錢包資訊
	wallet, err := getMemberWallet(rdb, memberId)
	if err != nil {
		return err
	}

	// 事務開始
	_, err = rdb.Do("MULTI")
	if err != nil {
		return err
	}

	// 大錢包減少金額
	if wallet.MainBalance < amount {
		rdb.Do("DISCARD")
		return err
	}
	wallet.MainBalance -= amount
	// 小錢包增加金額
	wallet.GameBalances[gameId] += amount
	err = setMemberWallet(rdb, memberId, wallet)
	if err != nil {
		rdb.Do("DISCARD")
		return err
	}

	// 執行事務
	_, err = rdb.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// 小錢包到主錢包轉帳
func transferSubToMain(rdb redis.Conn, memberId, gameId int, amount float64) error {
	// 會員錢包資訊
	wallet, err := getMemberWallet(rdb, memberId)
	if err != nil {
		return err
	}

	// 事務開始
	_, err = rdb.Do("MULTI")
	if err != nil {
		return err
	}

	// 小錢包減少金額
	if wallet.GameBalances[gameId] < amount {
		return err
	}
	wallet.GameBalances[gameId] -= amount
	// 大錢包增加金額
	wallet.MainBalance += amount

	err = setMemberWallet(rdb, memberId, wallet)
	if err != nil {
		rdb.Do("DISCARD")
		return err
	}

	// 執行事務
	_, err = rdb.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 主錢包點數吸入小錢包
func transferBringIn(rdb redis.Conn, memberId, gameId int) error {
	// 會員錢包資訊
	wallet, err := getMemberWallet(rdb, memberId)
	if err != nil {
		return err
	}

	_, err = rdb.Do("MULTI")
	if err != nil {
		return err
	}
	amount := wallet.MainBalance
	wallet.GameBalances[gameId] += amount
	wallet.MainBalance = 0

	err = setMemberWallet(rdb, memberId, wallet)
	if err != nil {
		rdb.Do("DISCARD")
		return err
	}

	// 執行事務
	_, err = rdb.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 全點數吸入主錢包
func transferBringOut(rdb redis.Conn, memberId int) error {
	// 會員錢包資訊
	wallet, err := getMemberWallet(rdb, memberId)
	if err != nil {
		return err
	}

	_, err = rdb.Do("MULTI")
	if err != nil {
		return err
	}

	for gameId, amount := range wallet.GameBalances {
		wallet.MainBalance += amount
		wallet.GameBalances[gameId] = 0
	}

	err = setMemberWallet(rdb, memberId, wallet)
	if err != nil {
		rdb.Do("DISCARD")
		return err
	}

	// 執行事務
	_, err = rdb.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 交易
func transferTransection(rdb redis.Conn, memberId int, amount float64, transferType string) error {
	// 會員錢包資訊
	wallet, err := getMemberWallet(rdb, memberId)
	if err != nil {
		return err
	}

	_, err = rdb.Do("MULTI")
	if err != nil {
		return err
	}

	if wallet.MainBalance < amount && transferType == "withdraw" {
		rdb.Do("DISCARD")
		return err
	}

	wallet.MainBalance -= amount

	err = setMemberWallet(rdb, memberId, wallet)
	if err != nil {
		rdb.Do("DISCARD")
		return err
	}

	_, err = rdb.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 遊戲帳務交易
func processGameTransaction(rdb redis.Conn, memberId, gameId int, amount float64) error {
	// 會員錢包資訊
	wallet, err := getMemberWallet(rdb, memberId)
	if err != nil {
		return err
	}

	// 事務開始
	_, err = rdb.Do("MULTI")
	if err != nil {
		return err
	}

	// 小錢包增加金額
	wallet.GameBalances[gameId] += amount
	err = setMemberWallet(rdb, memberId, wallet)
	if err != nil {
		rdb.Do("DISCARD")
		return err
	}

	// 執行事務
	_, err = rdb.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

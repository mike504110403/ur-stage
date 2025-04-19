package database

import (
	"github.com/mike504110403/goutils/dbconn"
)

// 以下為連線設定範例 需調整
// 連線字符定義在.env
const Envkey string = "MYSQL_URL"

// mysql使用的參數設定
const WALLET dbconn.DBName = "WALLET"
const GAME dbconn.DBName = "GAME"
const MEMBER dbconn.DBName = "MEMBER"
const SETTING dbconn.DBName = "Setting"

// 組裝用字串
const Wallet_dsn string = "Wallet"
const Setting_dsn string = "Setting"
const Member_dsn string = "Member"
const Game_dsn string = "Game"

var DB_Name_Map = map[dbconn.DBName]string{
	WALLET:  Wallet_dsn,
	SETTING: Setting_dsn,
	MEMBER:  Member_dsn,
	GAME:    Game_dsn,
}

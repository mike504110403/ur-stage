package database

import (
	"github.com/mike504110403/goutils/dbconn"
)

// 以下為連線設定範例 需調整
// 連線字符定義在.env
const Envkey string = "MYSQL_URL"

// mysql使用的參數設定
const WALLET dbconn.DBName = "WALLET"
const MEMBER dbconn.DBName = "MEMBER"
const SETTING dbconn.DBName = "Setting"
const ORDER dbconn.DBName = "Order"
const POINT dbconn.DBName = "Point"

// 組裝用字串
const Wallet_dsn string = "Wallet"
const Setting_dsn string = "Setting"
const Member_dsn string = "Member"
const Order_dsn string = "Order"
const Point_dsn string = "Point"

var DB_Name_Map = map[dbconn.DBName]string{
	WALLET:  Wallet_dsn,
	MEMBER:  Member_dsn,
	ORDER:   Order_dsn,
	SETTING: Setting_dsn,
	POINT:   Point_dsn,
}

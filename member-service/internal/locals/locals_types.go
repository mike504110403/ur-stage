package locals

// Key : 使用於locals的key
type LocalsKey = string

// KeyUserID : userID使用的KeyNames
const KeyUserID LocalsKey = "userID"

// KeyUserInfo : userInfo使用的KeyNames
const KeyUserInfo LocalsKey = "userInfo"

// KeyJWTToken : 於存放jwt的token的key name
const KeyJWTToken LocalsKey = "jwtToken"

type UserInfo struct {
	MemberId int
	UserName string
}

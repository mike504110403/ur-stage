package gb_elec

import (
	"fmt"
	"game_service/internal/service"

	mlog "github.com/mike504110403/goutils/log"
)

// TODO: 這邊才是對“遊戲”介面的實作 層次要抽開

// Apollo apollo遊戲商
type GBELECGamer struct {
}

// 進入遊戲
func (g *GBELECGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	res := ""

	url, err := lobbyLoginAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return "加入遊戲失敗", err
	} else {
		res = url
	}
	return res, nil
}

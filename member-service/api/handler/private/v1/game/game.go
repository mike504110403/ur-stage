package game

import (
	"database/sql"
	"errors"
	"member_service/internal/database"
	"sort"

	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

func getGameList(req GameListReq) ([]GameAgentList, error) {
	list := []GameAgentList{}
	db, err := database.GAME.DB()
	if err != nil {
		return list, err
	}
	donateType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "live_donate",
	}
	donateInt, err := donateType.Get()
	if err != nil {
		return list, err
	}

	queryStr := `
		SELECT id, agent_name, agent_type
		FROM GameAgent
		WHERE is_enable = 1
	`
	var rows *sql.Rows
	if req.GameType != nil {
		queryStr += " AND agent_type = ?"

		// 遊戲種類
		gametype := typeparam.TypeParam{
			MainType: "game_type",
			SubType:  typeparam.SubType(*req.GameType),
		}
		state, err := gametype.Get()
		if err != nil {
			return list, err
		}
		if state == 0 {
			return list, errors.New("game type not found")
		}
		rows, err = db.Query(queryStr, state)
		if err != nil {
			return list, err
		}
	} else {
		rows, err = db.Query(queryStr)
		if err != nil {
			return list, err
		}
	}

	defer rows.Close()
	for rows.Next() {
		agent := GameAgentList{}
		if err := rows.Scan(&agent.AgentId, &agent.AgentName, &agent.AgentType); err != nil {
			return list, err
		} else {
			if agent.AgentType != donateInt {
				list = append(list, agent)
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].AgentId > list[j].AgentId
	})
	return list, nil
}

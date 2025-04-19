package game

type (
	GameAgentList struct {
		AgentId   int    `json:"agent_id"`
		AgentName string `json:"agent_name"`
		AgentType int    `json:"agent_type"`
	}
	GameListReq struct {
		GameType *string `json:"game_type"`
	}
)

type (
	JoinGameRes struct {
		Url string `json:"url"`
	}
)

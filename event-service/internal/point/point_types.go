package point

import "time"

// 流水資訊
type BetFlow struct {
	Id        int       `db:"id"`
	MemberId  int       `db:"member_id"`
	BetFlow   float64   `db:"bet_flow"`
	TimeStart time.Time `db:"time_start"`
	TimeEnd   time.Time `db:"time_end"`
}

type MemberEventDto struct {
	ID           int     `db:"id"`
	BonusEventID int     `db:"event_id"`
	GateFlow     float64 `db:"gate_flow"`
}

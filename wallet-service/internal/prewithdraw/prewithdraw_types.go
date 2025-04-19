package prewithdraw

import event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"

type WalletGateDto struct {
	MemberId              int      `db:"member_id"`
	Amount                float64  `db:"balance"`
	WithdrawGate          *float64 `db:"withdraw_gate"`
	TihsdayWithdrawTimes  *int     `db:"thisday_withdraw_times"`
	ThisdayWithdrawAmount *float64 `db:"thisday_withdraw_amount"`
	AccBonus              *float64 `db:"acc_bonus"`
	WithdrawableBonus     *float64 `db:"withdrawable_bonus"`
}

// EventGateDto 事件提款門檻
type EventGateDto struct {
	MemberId              int      `db:"member_id"`
	Amount                float64  `db:"balance"`
	TihsdayWithdrawTimes  *int     `db:"thisday_withdraw_times"`
	ThisdayWithdrawAmount *float64 `db:"thisday_withdraw_amount"`
	Events                []event.MemberEvents
}

package point

import (
	"database/sql"
	"event_service/internal/database"
	"fmt"
	"time"

	mlog "github.com/mike504110403/goutils/log"
)

// 更新會員活動流水
func assignFlow(flow BetFlow) error {
	betFlow := flow.BetFlow

	db, err := database.POINT.DB()
	if err != nil {
		return err
	}
	queryStr := `
		SELECT id, event_id, gate_flow
		FROM MemberEvents
		WHERE member_id = ? AND is_release = 0
		ORDER BY create_at ASC
	`
	rows, err := db.Query(queryStr, flow.MemberId)
	if err != nil {
		return err
	}

	defer rows.Close()

	// 依時序分配流水
	memberEvents := map[int]MemberEventDto{}
	for rows.Next() {
		memberEvent := MemberEventDto{}
		if err := rows.Scan(
			&memberEvent.ID,
			&memberEvent.BonusEventID,
			&memberEvent.GateFlow,
		); err != nil {
			return err
		}

		memberEvent.GateFlow += betFlow
		if memberEvent.GateFlow > 0 {
			betFlow += memberEvent.GateFlow
			memberEvent.GateFlow = 0
			memberEvents[memberEvent.ID] = memberEvent
			continue
		} else {
			memberEvents[memberEvent.ID] = memberEvent
			break
		}
	}

	return updateMemberEvent(memberEvents, flow)
}

// 更新會員活動流水
func updateMemberEvent(memberEvents map[int]MemberEventDto, flow BetFlow) error {
	tx, err := database.POINT.TX()
	if err != nil {
		return err
	}

	for _, memberEvent := range memberEvents {
		if memberEvent.GateFlow == 0 {
			updateStr := `
				UPDATE MemberEvents
				SET gate_flow = 0, is_release = 1
				WHERE id = ?
			`
			_, err = tx.Exec(updateStr, memberEvent.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			updateStr := `
				UPDATE MemberEvents
				SET gate_flow = ?
				WHERE id = ?
			`
			_, err = tx.Exec(updateStr, memberEvent.GateFlow, memberEvent.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := insertRecognized(tx, memberEvent.ID, flow.Id, -memberEvent.GateFlow); err != nil {
			tx.Rollback()
			mlog.Error(fmt.Sprintf("流水認列失敗: %v", err))
			return err
		}

	}

	return tx.Commit()
}

// 寫入流水認列
func insertRecognized(tx *sql.Tx, eventId int, flowId int, amount float64) error {
	insertStr := `
		INSERT INTO RecognizeBetFlow (member_events_id, bet_flow_id, amount)
		VAlUES (?, ?, ?)
	`
	_, err := tx.Exec(insertStr, eventId, flowId, amount)
	if err != nil {
		return err
	}
	return nil
}

// 取得活動流水
func fetchBetFlows() ([]BetFlow, error) {
	db, err := database.POINT.DB()
	if err != nil {
		return nil, err
	}

	queryStr := `
		SELECT
			bf.id,
			bf.member_id,
			bf.bet_flow,
			bf.time_start,
			bf.time_end
		FROM
			Member.BetFlow_Minute AS bf
		LEFT JOIN
			RecognizeBetFlow AS rbf ON bf.id = rbf.bet_flow_id
		WHERE
			rbf.bet_flow_id IS NULL AND bf.time_end < ?
	`
	rows, err := db.Query(queryStr, time.Now().Add(-10*time.Minute))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	betFlows := []BetFlow{}
	for rows.Next() {
		betFlow := BetFlow{}
		if err := rows.Scan(
			&betFlow.Id,
			&betFlow.MemberId,
			&betFlow.BetFlow,
			&betFlow.TimeStart,
			&betFlow.TimeEnd,
		); err != nil {
			return nil, err
		}

		betFlows = append(betFlows, betFlow)
	}

	return betFlows, nil
}

// 取得活動流水筆數
func getBetFlowCounts() (int, error) {
	db, err := database.POINT.DB()
	if err != nil {
		return 0, err
	}

	queryStr := `
		SELECT COUNT(*)
		FROM Member.BetFlow_Minute AS bf
		LEFT JOIN RecognizeBetFlow AS rbf ON bf.id = rbf.bet_flow_id
		WHERE rbf.bet_flow_id IS NULL AND bf.time_end < ?
	`
	var count int
	if err := db.QueryRow(queryStr, time.Now().Add(-10*time.Minute)).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

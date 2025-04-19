package point

import (
	"database/sql"
	"event_service/internal/worker"
	"fmt"
	"sync"
	"time"

	mlog "github.com/mike504110403/goutils/log"
)

// CronGetBetFlow : 定時取得流水資訊
func CronGetBetFlow() {
	for {
		if betFlowCounts, err := getBetFlowCounts(); err != nil {
			mlog.Error(fmt.Sprintf("getBetFlowCounts error: %v", err))
			time.Sleep(10 * time.Minute)
			continue
		} else if betFlowCounts == 0 {
			time.Sleep(10 * time.Minute)
			continue
		} else {
			numWorkers := betFlowCounts / 100
			if numWorkers < 1 {
				numWorkers = 1
			}
			// 建立 worker
			jobs := make(chan interface{}, betFlowCounts)
			var wg sync.WaitGroup

			for w := 1; w <= numWorkers; w++ {
				wg.Add(1)
				go worker.Worker(w, jobs, ProcessBetFlow, &wg)
			}

			betFlows, err := fetchBetFlows()
			if err != nil {
				if err == sql.ErrNoRows {
					mlog.Info("no betFlows")
					time.Sleep(10 * time.Minute)
					continue
				}
				mlog.Error(fmt.Sprintf("fetchBetFlows error: %v", err))
				continue
			}
			for _, flow := range betFlows {
				jobs <- flow
			}
			// 關閉 betFlowJobs channel，等待所有 worker 完成
			close(jobs)
			wg.Wait()

			time.Sleep(10 * time.Minute)
		}
	}
}

// ProcessBetFlow : 處理流水資訊
func ProcessBetFlow(f interface{}) {
	flow, ok := f.(BetFlow)
	if !ok {
		mlog.Error(fmt.Sprintf("ProcessBetFlow error: %v", "type assertion error"))
		return
	}
	if err := assignFlow(flow); err != nil {
		mlog.Error(fmt.Sprintf("assignFlow error: %v", err))
	}
}

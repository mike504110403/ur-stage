package worker

import (
	"sync"

	mlog "github.com/mike504110403/goutils/log"
)

// worker
func Worker(workerNumber int, jobs <-chan interface{}, process func(interface{}), wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		mlog.Info("活動流水計算 start")
		process(j)
	}
}

package framework

import (
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	MRand "math/rand"

	RB "CuttleTest/robot"
)

var (
	curStartTime, nextStartTime int64
)

type RunCase func(robot RB.RobotIF, stepId int)

type RunCaseInfo struct {
	RunCase
	StepCount int
}

//-------------------------------------------------

func FanInRobot(createFunc RB.CreateRobotFunc) chan RB.RobotIF {

	curStartTime = time.Now().UnixNano() / 1000000
	nextStartTime = curStartTime + (int64)(1000*PkgInterval)
	robots := make(chan RB.RobotIF, 2000)

	for i := 0; i < RobotCount; i++ {
		go func() {
			var robot RB.RobotIF

			createFunc(&robot, i, 1)

			robots <- robot
		}()
	}

	return robots
}

func DoTest(robots chan RB.RobotIF, runCaseArr []RunCaseInfo) {
	MRand.Seed(int64(time.Now().Nanosecond()))
	reqSem := make(chan struct{}, ReqCount) //限制同一时间发出的请求数

	var i int64
	for r := range robots {
		if r.GetNextStartTime() <= (time.Now().UnixNano() / 1000000) {
			reqSem <- struct{}{}
			go func() {
				if sceneID >= len(runCaseArr) {
					log.Printf("sceneId: %v, error", sceneID)
					return
				}

				runCaseInfo := runCaseArr[sceneID]
				//runCaseInfo.RunCase(rIF, r.NextStepID)
				runCaseInfo.RunCase(r, r.GetNextStepID())

				r.SetNextStepID((r.GetNextStepID() + 1) % runCaseInfo.StepCount)
				//r.NextStepID = (r.NextStepID + 1) % runCaseInfo.StepCount
				//r.NextStartTime = nextStartTime + (int64)(MRand.New(MRand.NewSource(time.Now().UnixNano())).Intn(1000*PkgInterval))
				startTime := r.GetNextStartTime() + (int64)(MRand.New(MRand.NewSource(time.Now().UnixNano())).Intn(1000*PkgInterval))
				r.SetNextStartTime(startTime)

				if atomic.AddInt64(&i, 1)%2000 == 0 {
					debug.FreeOSMemory()
				}

				<-reqSem

				r.SetCases(nil)
				r.SetErr(nil)

				robots <- r
			}()
		} else {
			robots <- r
		}

		if (time.Now().UnixNano() / 1000000) > nextStartTime {
			lock := &sync.Mutex{}
			lock.Lock()

			curStartTime = nextStartTime
			nextStartTime = curStartTime + (int64)(1000*PkgInterval)

			lock.Unlock()
		}

	}
}

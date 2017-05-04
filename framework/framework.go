package framework

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	MRand "math/rand"

	RB "github.com/homike/cuttletest/robot"
)

var (
	curStartTime, nextStartTime int64
)

type RunCase func(robot *RB.Robot, stepId int)

type RunCaseInfo struct {
	RunCase
	StepCount int
}

type InitCase func(r *RB.Robot, stepIndex int)

func FanInRobot(initCaseFunc InitCase) chan *RB.Robot {

	curStartTime = time.Now().UnixNano() / 1000000
	nextStartTime = curStartTime + (int64)(1000*PkgInterval)
	robots := make(chan *RB.Robot, 2000)

	for i := 0; i < RobotCount; i++ {

		name := ""
		if i == 0 {
			name = RobotName
		} else {
			name = fmt.Sprintf("%v%v", RobotName, i)
		}
		go func() {

			robot := &RB.Robot{
				RobotIndex: i,
				Name:       name,
				Password:   "123456",
			}
			initCaseFunc(robot, 1)

			robot.NextStepID = 0 //maRand.New(maRand.NewSource(time.Now().UnixNano())).Intn(totalStepNum[sceneId])
			robot.NextStartTime = time.Now().UnixNano()/1000000 + (int64)(MRand.New(MRand.NewSource(time.Now().UnixNano())).Intn(1000*PkgInterval))
			robots <- robot
		}()
	}

	return robots
}

func DoTest(robots chan *RB.Robot, runCaseArr []RunCaseInfo) {
	MRand.Seed(int64(time.Now().Nanosecond()))
	reqSem := make(chan struct{}, ReqCount) //限制同一时间发出的请求数

	var i int64
	for r := range robots {
		r := r
		if r.NextStartTime <= (time.Now().UnixNano() / 1000000) {
			reqSem <- struct{}{}
			go func() {
				if sceneID >= len(runCaseArr) {
					log.Printf("sceneId: %v, error", sceneID)
					return
				}

				runCaseInfo := runCaseArr[sceneID]
				runCaseInfo.RunCase(r, r.NextStepID)

				r.NextStepID = (r.NextStepID + 1) % runCaseInfo.StepCount
				r.NextStartTime = nextStartTime + (int64)(MRand.New(MRand.NewSource(time.Now().UnixNano())).Intn(1000*PkgInterval))

				if atomic.AddInt64(&i, 1)%2000 == 0 {
					debug.FreeOSMemory()
				}

				<-reqSem

				r.Cases = nil
				r.Err = nil

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

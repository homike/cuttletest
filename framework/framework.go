package framework

import (
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	MRand "math/rand"

	RB "cuttletest/robot"
)

var (
	curStartTime, nextStartTime int64
)

type RunCase func(robot *RB.Robot, stepId int)

type RunCaseInfo struct {
	RunCase
	StepCount int
}

type InitCase func(r *RB.Robot, robotIndex, sceneID int)

func FanInRobot(initCaseFunc InitCase) chan *RB.Robot {

	curStartTime = time.Now().UnixNano() / 1000000
	nextStartTime = curStartTime + (int64)(1000*RobotCfg.PkgInterval)
	robots := make(chan *RB.Robot, 2000)

	for i := 0; i < RobotCfg.RobotCount; i++ {
		go func() {

			robot := &RB.Robot{
				RobotIndex: i,
			}
			initCaseFunc(robot, i, RobotCfg.SceneID)

			robot.NextStepID = 0 //maRand.New(maRand.NewSource(time.Now().UnixNano())).Intn(totalStepNum[sceneId])
			robot.NextStartTime = time.Now().UnixNano()/1000000 + (int64)(MRand.New(MRand.NewSource(time.Now().UnixNano())).Intn(1000*RobotCfg.PkgInterval))
			robots <- robot
		}()
	}

	return robots
}

func DoTest(robots chan *RB.Robot, runCaseArr []RunCaseInfo) {
	MRand.Seed(int64(time.Now().Nanosecond()))
	reqSem := make(chan struct{}, RobotCfg.ReqCount) //限制同一时间发出的请求数

	var i int64
	for r := range robots {
		r := r
		if r.NextStartTime <= (time.Now().UnixNano() / 1000000) {
			reqSem <- struct{}{}
			go func() {
				if RobotCfg.SceneID >= len(runCaseArr) {
					log.Printf("sceneId: %v, error", RobotCfg.SceneID)
					return
				}

				runCaseInfo := runCaseArr[RobotCfg.SceneID]
				runCaseInfo.RunCase(r, r.NextStepID)

				r.NextStepID = (r.NextStepID + 1) % runCaseInfo.StepCount
				r.NextStartTime = nextStartTime + (int64)(MRand.New(MRand.NewSource(time.Now().UnixNano())).Intn(1000*RobotCfg.PkgInterval))

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
			nextStartTime = curStartTime + (int64)(1000*RobotCfg.PkgInterval)

			lock.Unlock()
		}
	}
}

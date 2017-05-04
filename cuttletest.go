package cuttletest

import (
	"github.com/homike/cuttletest/framework"
)

func Run(configPath string, initCaseFunc framework.InitCase, runCaseList []framework.RunCaseInfo) {
	framework.InitConfig(configPath)

	robots := framework.FanInRobot(initCaseFunc)
	framework.DoTest(robots, runCaseList)
}

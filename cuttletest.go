package cuttletest

import (
	"CuttleTest/framework"
)

// Run Test
func Run(configPath string, initCaseFunc framework.InitCase, runCaseList []framework.RunCaseInfo) {
	framework.InitConfig(configPath)

	robots := framework.FanInRobot(initCaseFunc)
	framework.DoTest(robots, runCaseList)
}

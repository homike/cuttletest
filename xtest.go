package xtest

import (
	RB "github.com/homike/cuttletest/robot"
)

func main() {
	robots := RB.FanInRobot(50)
	RB.DoTest(robots)
}

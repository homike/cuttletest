package robot

import (
	"fmt"

	CS "github.com/homike/cuttletest/cases"
)

type Robot struct {
	RobotIndex int
	Name       string
	Password   string
	Cases      []CS.Case
}

func (r *Robot) initCase() {
}

// AddCase Add Test Case
func (r *Robot) AddCase(c CS.Case) {
	r.Cases = append(r.Cases, c)
}

func (r *Robot) Play() {
	for _, c := range r.Cases {
		c.Assemble()
		c.Do()
	}
}

func FanInRobot(robotCount int) chan *Robot {

	robots := make(chan *Robot, 20000)

	for i := 0; i < robotCount; i++ {

		name := fmt.Sprintf("robot%v", i)
		go func() {
			robot := &Robot{
				RobotIndex: i,
				Name:       name,
				Password:   "123456",
			}

			robot.initCase()
			robots <- robot
		}()
	}

	return robots
}

func DoTest(robots chan *Robot) {

	for robot := range robots {
		go func() {
			robots <- robot
		}()
	}

}

package robot

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"

	CS "CuttleTest/cases"
)

var client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

type CreateRobotFunc func(r *RB.RobotIF, robotIndex int, stepIndex int)

type RobotIF interface {
	AddCase(c CS.CaseFunction)
	Act(i int, m map[string]string)
	Play(m map[string]string) error

	// RobotIndex
	GetRobotIndex() int
	SetRobotIndex(index int)
	// StepID
	SetNextStepID(stepID int)
	GetNextStepID() int
	// StartTime
	SetNextStartTime(nextTime int64)
	GetNextStartTime() int64
	// Cases
	SetCases([]CS.Case)
	// ERR
	SetErr(error)
}

type Robot struct {
	RobotIndex    int
	Cases         []CS.Case
	Err           error
	NextStepID    int
	NextStartTime int64
}

func (r *Robot) GetRobotIndex() int {
	return r.RobotIndex
}

func (r *Robot) SetRobotIndex(index int) {
	r.RobotIndex = index
}

func (r *Robot) GetNextStepID() int {
	return r.NextStepID
}

func (r *Robot) SetNextStepID(stepID int) {
	r.NextStepID = stepID
}

func (r *Robot) SetNextStartTime(nextTime int64) {
	r.NextStartTime = nextTime
}

func (r *Robot) GetNextStartTime() int64 {
	return r.NextStartTime
}

func (r *Robot) SetCases(c []CS.Case) {
	r.Cases = c
}

func (r *Robot) SetErr(err error) {
	r.Err = err
}

func (r *Robot) AddCase(c CS.CaseFunction) {
	r.Cases = append(r.Cases, CS.Case{CaseFunc: c})
}

func (r *Robot) Act(i int, m map[string]string) {
	if r.Err != nil {
		return
	}

	info := make(url.Values)

	for k, v := range m {
		info.Add(k, v)
	}

	r.Cases[i].CaseFunc.Assemble(info)
	r.Cases[i].CaseRet, r.Err = r.Cases[i].CaseFunc.Do(client)
	if r.Err == nil {
		return
	}
}

func (r *Robot) Play(m map[string]string) error {
	for i, _ := range r.Cases {

		r.Act(i, m)

		if r.Err != nil {
			errStr := r.Err.Error()
			index := strings.Index(errStr, "query")
			if index > 0 {
				errStr = errStr[0:index]
			}
			log.Println(errStr)
			return r.Err
		}
	}

	return nil
}

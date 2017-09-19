package robot

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"

	"time"

	"CuttleTest/mode"
)

var client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

var (
	TotalReqCount, TotalReqTime int64
)

type UserRoboter interface {
	ActExtraFunc() map[string]string
}

type Robot struct {
	RobotIndex    int
	Cases         []mode.Case
	Err           error
	NextStepID    int
	NextStartTime int64

	// Extra Data
	ExtraData UserRoboter
}

func (r *Robot) AddCase(c mode.CaseFunction) {
	r.Cases = append(r.Cases, mode.Case{CaseFunc: c})
}

func (r *Robot) act(i int, m map[string]string) {
	if r.Err != nil {
		return
	}

	info := make(url.Values)

	extraParas := r.ExtraData.ActExtraFunc()
	// if r.AccountID != "" && r.Token != "" {
	// 	info.Add("account_id", r.AccountID)
	// 	info.Add("user_token", r.Token)

	// 	for k, v := range m {
	// 		info.Add(k, v)
	// 	}
	// }

	for k, v := range extraParas {
		info.Add(k, v)
	}

	for k, v := range m {
		info.Add(k, v)
	}

	r.Cases[i].CaseFunc.Assemble(info)
	r.Cases[i].CaseRet, r.Err = r.Cases[i].CaseFunc.Do(client)
	if r.Err == nil {
		return
	}
}

func (r *Robot) Play() error {
	for i, _ := range r.Cases {

		t1 := time.Now().UnixNano()
		r.act(i, nil)

		t2 := time.Now().UnixNano()

		TotalReqCount = TotalReqCount + 1
		TotalReqTime = TotalReqTime + (t2 - t1)

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

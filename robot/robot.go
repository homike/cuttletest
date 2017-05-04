package robot

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"

	"time"

	CS "github.com/homike/cuttletest/cases"
)

var client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

var (
	TotalReqCount, TotalReqTime int64
)

type Robot struct {
	RobotIndex    int
	Cases         []CS.Case
	Err           error
	NextStepID    int
	NextStartTime int64
	//extend
	AccountID string
	Token     string
	Name      string
	Password  string
}

func (r *Robot) AddCase(c CS.CaseFunction) {
	r.Cases = append(r.Cases, CS.Case{CaseFunc: c})
}

func (r *Robot) act(i int, m map[string]string) {
	if r.Err != nil {
		return
	}

	info := make(url.Values)

	if r.AccountID != "" && r.Token != "" {
		info.Add("account_id", r.AccountID)
		info.Add("user_token", r.Token)

		for k, v := range m {
			info.Add(k, v)
		}
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

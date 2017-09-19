package mode

import (
	"net/http"
	"net/url"
)

type Result map[string]interface{}

type CaseFunction interface {
	Assemble(form url.Values)
	Do(client *http.Client) (Result, error)
}

type Case struct {
	CaseFunc CaseFunction
	CaseRet  Result
}

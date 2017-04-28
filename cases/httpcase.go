package cases

import (
	"net/http"
	"net/url"
)

type Result map[string]interface{}

type Case interface {
	Assemble(form url.Values)
	Do(client *http.Client) (Result, error)
}

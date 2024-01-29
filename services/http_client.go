package services

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	DefaultHTTPClient HTTPClient = &http.Client{}
)

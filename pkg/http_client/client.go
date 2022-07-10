package http_client

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func Init() {
	Client = &http.Client{}
}

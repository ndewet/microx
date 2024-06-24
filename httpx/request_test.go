package httpx

import (
	"net/http"
	"net/url"
)

func CreateMockHTTPRequest(method Method, path string) *http.Request {
	return &http.Request{
		Method: string(method),
		URL:    &url.URL{Path: path},
	}
}

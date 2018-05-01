package interfaces

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OfficialHttpClientAdapter struct {
	client *http.Client
}

func (c *OfficialHttpClientAdapter) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

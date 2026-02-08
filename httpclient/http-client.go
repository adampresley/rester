package httpclient

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type mockHttpClient struct {
	resp *http.Response
	err  error
}

func NewMockHttpClient(resp *http.Response, err error) HttpClient {
	return &mockHttpClient{
		resp: resp,
		err:  err,
	}
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

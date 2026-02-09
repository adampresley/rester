package clientoptions

import (
	"encoding/base64"
	"net/http"

	"github.com/adampresley/rester/contenttype"
	"github.com/adampresley/rester/httpclient"
)

type ClientOptions struct {
	BaseURL                   string
	CustomContentTypeHandlers map[string]contenttype.ContentTypeHandler
	Debug                     bool
	Headers                   map[string]string
	BasicAuthHeader           string
	HttpClient                httpclient.HttpClient
}

type ClientOption func(*ClientOptions)

func New(baseURL string, options ...ClientOption) *ClientOptions {
	result := &ClientOptions{
		BaseURL:                   baseURL,
		CustomContentTypeHandlers: map[string]contenttype.ContentTypeHandler{},
		HttpClient:                http.DefaultClient,
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func WithBasicAuth(username, password string) ClientOption {
	return func(s *ClientOptions) {
		base64Encoded := base64.StdEncoding.EncodeToString(
			[]byte(username + ":" + password),
		)

		s.BasicAuthHeader = base64Encoded
	}
}

func WithCustomContentTypeHandler(contentType string, handler contenttype.ContentTypeHandler) ClientOption {
	return func(s *ClientOptions) {
		s.CustomContentTypeHandlers[contentType] = handler
	}
}

func WithDebug(debug bool) ClientOption {
	return func(s *ClientOptions) {
		s.Debug = debug
	}
}

func WithHeaders(headers map[string]string) ClientOption {
	return func(s *ClientOptions) {
		s.Headers = headers
	}
}

func WithHttpClient(client httpclient.HttpClient) ClientOption {
	return func(s *ClientOptions) {
		s.HttpClient = client
	}
}

package testing

import (
	"bytes"
	"easycoding/pkg/http_client"
	"net/http"
)

type Dofunc func(*http.Request) (*http.Response, error)

type HTTPMockClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m *HTTPMockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func NewHttpMockClient(f Dofunc) *HTTPMockClient {
	return &HTTPMockClient{f}
}

type BodyReader struct {
	*bytes.Reader
}

func (*BodyReader) Close() error {
	return nil
}

func NewBodyReader(b []byte) *BodyReader {
	return &BodyReader{bytes.NewReader(b)}
}

func SetupHttpClient(f Dofunc) {
	http_client.Client = NewHttpMockClient(f)
}

package jotform

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy Response"))}, nil
}

func NewTestClient(mockHttp *MockHttpClient) *jotformAPIClient {
	client := NewJotFormAPIClient("api-key", "json", false)
	client.HttpClient = mockHttp
	return client
}

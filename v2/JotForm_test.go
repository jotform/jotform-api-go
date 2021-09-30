package jotform_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	jotform "github.com/jotform/jotform-api-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	t.Run("happy - uses default URL", func(t *testing.T) {
		submissionID := int64(123)

		var reqURL string
		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				reqURL = req.URL.String()
				return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy Response"))}, nil
			}},
		)

		_, _ = client.GetSubmission(submissionID)
		assert.Equal(t, reqURL, fmt.Sprintf("https://api.jotform.com/v1/user/submission/%d", submissionID))
	})

	t.Run("happy - uses enterprise URL", func(t *testing.T) {
		submissionID := int64(456)
		enterpriseURL := "https://enterprise.jotform.com/API"

		var reqURL string
		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				reqURL = req.URL.String()
				return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy Response"))}, nil
			}},
		)
		client.BaseURL = enterpriseURL

		_, _ = client.GetSubmission(submissionID)
		assert.Equal(t, reqURL, fmt.Sprintf("%s/v1/user/submission/%d", enterpriseURL, submissionID))
	})
}

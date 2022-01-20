package jotform_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	jotform "github.com/jotform/jotform-api-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestDownloadRichPDFSubmission(t *testing.T) {
	t.Run("happy - returns byte response", func(t *testing.T) {
		pdfString := "Pretend this is a PDF"
		formID := "123"
		submissionID := "456"

		var reqURL string

		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				reqURL = req.URL.String()
				return &http.Response{
					Request:    req,
					StatusCode: 200,
					Status:     "200 OK",
					Body:       ioutil.NopCloser(bytes.NewBufferString(pdfString))}, nil
			}},
		)

		res, err := client.DownloadRichPDFSubmission(formID, submissionID)
		assert.Nil(t, err)
		assert.Equal(t, reqURL, fmt.Sprintf("https://api.jotform.com/v1/pdf-converter/%s/fill-pdf?submissionID=%s", formID, submissionID))
		assert.Equal(t, pdfString, string(res))
	})

	t.Run("sad - permission denied", func(t *testing.T) {
		authFailedBody := `{"responseCode":401,"message":"You're not authorized to use (\/pdf-converter-id-fill-pdf) ","content":"","duration":"108.22ms","info":"https:\/\/api.jotform.com\/docs#pdf-converter-id-fill-pdf"}`
		formID := "123"
		submissionID := "456"

		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Request:    req,
					StatusCode: 401,
					Status:     "401 Not Authorized",
					Body:       ioutil.NopCloser(bytes.NewBufferString(authFailedBody))}, nil
			}},
		)

		_, err := client.DownloadRichPDFSubmission(formID, submissionID)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "401")
	})

	t.Run("sad - form has no PDF", func(t *testing.T) {
		noPDFBody := `{"responseCode":400,"message":"draw-pdf-answers Request Failed","content":"","duration":"98.08ms","info":"https:\/\/api.jotform.com\/docs#pdf-converter-id-fill-pdf"}`
		formID := "123"
		submissionID := "456"

		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Request:    req,
					StatusCode: 400,
					Status:     "400 Bad Request",
					Body:       ioutil.NopCloser(bytes.NewBufferString(noPDFBody))}, nil
			}},
		)

		_, err := client.DownloadRichPDFSubmission(formID, submissionID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, jotform.ErrNotImplemented))
	})
}

func TestDownloadSimplePDFSubmission(t *testing.T) {

	t.Run("happy - returns byte response", func(t *testing.T) {
		pdfString := "Pretend this is a PDF"
		formID := "123"
		submissionID := "456"
		reportID := "789"

		var reqURL string

		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				reqURL = req.URL.String()
				return &http.Response{
					Request:    req,
					StatusCode: 200,
					Status:     "200 OK",
					Body:       ioutil.NopCloser(bytes.NewBufferString(pdfString))}, nil
			}},
		)

		res, err := client.DownloadSimplePDFSubmission(formID, submissionID, reportID)
		assert.Nil(t, err)
		assert.Equal(t, reqURL, fmt.Sprintf("https://api.jotform.com/v1/generatePDF?download=1&formid=%s&reportid=789&submissionid=%s", formID, submissionID))
		assert.Equal(t, pdfString, string(res))
	})

	t.Run("happy - no report ID", func(t *testing.T) {
		pdfString := "Pretend this is a PDF"
		formID := "123"
		submissionID := "456"
		reportID := ""

		var reqURL string

		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				reqURL = req.URL.String()
				return &http.Response{
					Request:    req,
					StatusCode: 200,
					Status:     "200 OK",
					Body:       ioutil.NopCloser(bytes.NewBufferString(pdfString))}, nil
			}},
		)

		res, err := client.DownloadSimplePDFSubmission(formID, submissionID, reportID)
		assert.Nil(t, err)
		assert.Equal(t, reqURL, fmt.Sprintf("https://api.jotform.com/v1/generatePDF?download=1&formid=%s&submissionid=%s", formID, submissionID))
		assert.Equal(t, pdfString, string(res))
	})

	t.Run("sad - permission denied", func(t *testing.T) {
		authFailedBody := `{"responseCode":401,"message":"Authorization error for user()-form(123)-token()!","content":"","duration":"82.04ms","info":"https:\/\/api.jotform.com\/docs#generatePDF"}`
		formID := "123"
		submissionID := "456"
		reportID := "789"

		client := jotform.NewTestClient(
			&jotform.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Request:    req,
					StatusCode: 401,
					Status:     "401 Not Authorized",
					Body:       ioutil.NopCloser(bytes.NewBufferString(authFailedBody))}, nil
			}},
		)

		_, err := client.DownloadSimplePDFSubmission(formID, submissionID, reportID)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "401")
	})
}

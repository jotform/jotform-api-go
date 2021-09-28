package jotform

import (
	"fmt"
	"io/ioutil"
)

// DownloadFormattedPDFSubmission returns a PDF
// for the provided submissionID and formID
// that was specifically formatted for that formID,
// ie. the form was created in Jotform from a PDF.
func (client jotformAPIClient) DownloadFormattedPDFSubmission(formID, submissionID string) ([]byte, error) {
	resp, err := client.HttpClient.Do(client.newRequest(
		fmt.Sprintf("pdf-converter/%s/fill-pdf", submissionID),
		map[string]string{
			"formid": formID,
		},
		"GET",
	))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Jotform API request for '%s' failed: %s", resp.Request.URL, resp.Status)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// DownloadSimplePDFSubmission returns a PDF
// with the form names and values filled out
// but not formatted in any special way
func (client jotformAPIClient) DownloadSimplePDFSubmission(formID, submissionID string) ([]byte, error) {
	resp, err := client.HttpClient.Do(client.newRequest(
		"generatePDF",
		map[string]string{
			"formid":       formID,
			"submissionid": submissionID,
			"download":     "1",
		},
		"GET",
	))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Jotform API request for '%s' failed: %s", resp.Request.URL, resp.Status)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

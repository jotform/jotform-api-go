package jotform

import (
	"errors"
	"fmt"
	"io/ioutil"
)

var ErrNotImplemented = errors.New("Not Implemented")

// DownloadRichPDFSubmission returns a PDF
// for the provided submissionID and formID
// that was specifically formatted for that formID,
// ie. the form was created in Jotform from a PDF.
func (client jotformAPIClient) DownloadRichPDFSubmission(formID, submissionID string) ([]byte, error) {
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
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		// This is a response like:
		// {"responseCode":400,"message":"draw-pdf-answers Request Failed","content":"","duration":"98.08ms","info":"https:\/\/api.jotform.com\/docs#pdf-converter-id-fill-pdf"}
		return nil, fmt.Errorf("Jotform form %s does not have an associated PDF: %w", formID, ErrNotImplemented)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Jotform API request for '%s' failed: %s", resp.Request.URL, resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// DownloadSimplePDFSubmission returns a PDF
// with the form names and values filled out
// formatted by the provided reportID.
// If no reportID is provided or the provided ID does not exist,
// this will default to the first PDF listed on the PDF Editor.
// If no PDFs exist on the PDF editor, this will generate one.
func (client jotformAPIClient) DownloadSimplePDFSubmission(formID, submissionID, reportID string) ([]byte, error) {
	query := map[string]string{
		"formid":       formID,
		"submissionid": submissionID,
		"download":     "1",
	}
	if reportID != "" {
		query["reportid"] = reportID
	}
	resp, err := client.HttpClient.Do(client.newRequest(
		"generatePDF",
		query,
		"GET",
	))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Jotform API request for '%s' failed: %s", resp.Request.URL, resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

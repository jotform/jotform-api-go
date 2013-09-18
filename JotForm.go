// JotForm API - GO Client
// copyright   2013 Interlogy, LLC.
// link        http://www.jotform.com
// version     1.0
// package     JotFormAPI

package jotform

import(
    "fmt"
    "strconv"
    "net/http"
    "net/url"   
    "io/ioutil"
    "os"
    "encoding/json"
    "encoding/xml"
    "strings"
    "bytes"
)

const baseURL = "http://api.jotform.com"
const apiVersion = "v1"

type jotformAPIClient struct{
    apiKey string
    outputType string
    debugMode bool
}

func NewJotFormAPIClient(apiKey string, outputType string, debugMode bool) *jotformAPIClient {
    client := &jotformAPIClient{apiKey, strings.ToLower(outputType), debugMode}

    return client
}

func (client jotformAPIClient) GetOutputType() string {return client.outputType}
func (client *jotformAPIClient) SetOutputType(value string) {client.outputType = value}

func (client jotformAPIClient) GetDebugMode() bool {return client.debugMode}
func (client *jotformAPIClient) SetDebugMode(value bool) {client.debugMode = value}

func (client jotformAPIClient) debug(str interface{}) {
    if client.debugMode {
        fmt.Println(str)
    }
}

func (client jotformAPIClient) executeHttpRequest(requestPath string, params interface{}, method string) []byte {

    if client.outputType != "json" {
        requestPath = requestPath + ".xml"
    }

    var path = baseURL + "/" + apiVersion + "/" + requestPath
    client.debug(path)

    var response *http.Response
    var request *http.Request
    var err error

    client.debug(params)

    if method == "GET" {
        if params != "" {
            data := params.(map[string]string)
            values := make(url.Values)

            for k, _ := range data {
                values.Set(k, data[k])
            }
            path = path +"?"+ values.Encode()
        }

        request, err = http.NewRequest("GET", path, nil)
        request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        request.Header.Add("apiKey", client.apiKey)
        response, err = http.DefaultClient.Do(request)
    } else if method == "POST" {
        data := params.(map[string]string)
        values := make(url.Values)

        for k, _ := range data {
            values.Set(k, data[k])
        }

        request, err = http.NewRequest("POST", path, strings.NewReader(values.Encode()))
        request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        request.Header.Add("apiKey", client.apiKey)
        response, err = http.DefaultClient.Do(request)
    } else  if method == "DELETE" {
        request, err = http.NewRequest("DELETE", path, nil)
        request.Header.Add("apiKey", client.apiKey)
        response, err = http.DefaultClient.Do(request)
    } else if method == "PUT" {
        parameters := params.([]byte)

        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        } else {
            request, err = http.NewRequest("PUT", path, bytes.NewBuffer(parameters))
            request.Header.Add("apiKey", client.apiKey)
            response, err = http.DefaultClient.Do(request)
        }
    }

    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }

        if client.outputType == "json" {
            var f interface{}
            json.Unmarshal(contents, &f)
            result := f.(map[string]interface{})["content"]
            content, err := json.Marshal(result)

            if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
            } else {
                return content   
            }
        } else if client.outputType == "xml" {
            var f interface{}
            xml.Unmarshal(contents, &f)
            return contents
        }
    }

    return nil
}

func createConditions(offset string, limit string, filter map[string]string, orderby string) map[string]string {
    args := map[string]interface{}{
        "offset": offset,
        "limit": limit,
        "filter": filter,
        "orderby": orderby,
    }

    params := make(map[string]string)

    for k, _ := range args {
        if k == "filter" {
            filterObj, err := json.Marshal(filter)

            if err == nil {
                params["filter"] = string(filterObj)
            }
        }else {
            params[k] = args[k].(string)
        }
    }
    return params
}

func createHistoryQuery (action string, date string, sortBy string, startDate string, endDate string) map[string]string {
    args := map[string]string {
        "action": action,
        "date": date,
        "sortBy": sortBy,
        "startDate": startDate,
        "endDate": endDate,
    }

    params := make(map[string]string)

    for k, _ := range args {
        if args[k] != "" {
            params[k] = args[k]
        }
    }
    return params
}

//GetUser
//Get user account details for a JotForm user.
//Returns user account type, avatar URL, name, email, website URL and account limits.
func (client jotformAPIClient) GetUser() []byte {
    return client.executeHttpRequest("user", "", "GET")
}

//GetUsage
//Get number of form submissions received this month
//Returns number of submissions, number of SSL form submissions, payment form submissions and upload space used by user.
func (client jotformAPIClient) GetUsage() []byte {
    return client.executeHttpRequest("user/usage", "", "GET")
}

//GetForms
//Get a list of forms for this account
//offset (string): Start of each result set for form list.
//limit (string): Number of results in each result set for form list.
//filter (map[string]string): Filters the query results to fetch a specific form range.
//orderBy (string): Order results by a form field name.
//Returns basic details such as title of the form, when it was created, number of new and total submissions.
func (client jotformAPIClient) GetForms(offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("user/forms", params, "GET")
}

//GetSubmissions
//Get a list of submissions for this account
//offset (string): Start of each result set for form list.
//limit (string): Number of results in each result set for form list.
//filter (map[string]string): Filters the query results to fetch a specific form range.
//orderBy (string): Order results by a form field name.
//Returns basic details such as title of the form, when it was created, number of new and total submissions.
func (client jotformAPIClient) GetSubmissions(offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("user/submissions", params, "GET")
}

//GetSubusers
//Get a list of sub users for this account
//Returns list of forms and form folders with access privileges.
func (client jotformAPIClient) GetSubusers() []byte {
    return client.executeHttpRequest("user/subusers", "", "GET")
}

//GetFolders
//Get a list of form folders for this account
//Returns name of the folder and owner of the folder for shared folders.
func (client jotformAPIClient) GetFolders() []byte {
    return client.executeHttpRequest("user/folders", "", "GET")
}

//GetReports
//List of URLS for reports in this account
//Returns reports for all of the forms. ie. Excel, CSV, printable charts, embeddable HTML tables.
func (client jotformAPIClient) GetReports() []byte {
    return client.executeHttpRequest("user/reports", "", "GET")
}

//Update user's settings
//New user setting values with setting keys
//Returns changes on user settings
func (client jotformAPIClient) GetSettings() []byte {
    return client.executeHttpRequest("user/settings", "", "GET")
}

//GetSettings
//Get user's settings for this account
//Returns user's time zone and language.
func (client jotformAPIClient) UpdateSettings(settings map[string]string) []byte {
    return client.executeHttpRequest("user/settings", settings, "POST")
}

//GetHistory
//Get user activity log
//action (string): Filter results by activity performed. Default is 'all'.
//date (string): Limit results by a date range. If you'd like to limit results by specific dates you can use startDate and endDate fields instead.
//sortBy (string): Lists results by ascending and descending order.
//startDate (string): Limit results to only after a specific date. Format: MM/DD/YYYY.
//endDate (string): Limit results to only before a specific date. Format: MM/DD/YYYY.
//Returns activity log about things like forms created/modified/deleted, account logins and other operations.
func (client jotformAPIClient) GetHistory(action string, date string, sortBy string, startDate string, endDate string) []byte {
    var params = createHistoryQuery(action, date, sortBy, startDate, endDate)

    return client.executeHttpRequest("user/history", params, "GET")
}

//GetForm
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns form ID, status, update and creation dates, submission count etc.
func (client jotformAPIClient) GetForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10), "", "GET")
}

//GetFormQuestions
//Get a list of all questions on a form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns question properties of a form.
func (client jotformAPIClient) GetFormQuestions(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", "", "GET")
}

//GetFormQuestion
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//qid (int): Identifier for each question on a form. You can get a list of question IDs from /form/{id}/questions.
//Returns question properties like required and validation.
func (client jotformAPIClient) GetFormQuestion(formID int64, qid int) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), "", "GET")
}

//GetFormSubmission
//List of a form submissions.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//offset (string): Start of each result set for form list.
//limit (string): Number of results in each result set for form list.
//filter (map[string]string): Filters the query results to fetch a specific form range.
//orderBy (string): Order results by a form field name.
//Returns submissions of a specific form.
func (client jotformAPIClient) GetFormSubmissions(formID int64, offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/submissions", params, "GET")
}

//CreateFormSubmission
//Submit data to this form using the API
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//submission (map[string]string): Submission data with question IDs.
//Returns posted submission ID and URL.
func (client jotformAPIClient) CreateFormSubmission(formId int64, submission map[string]string) []byte {
    data := make(map[string]string)

    for k, _ := range submission {
        if strings.Contains(k, "_") {
            data["submission[" + k[0:strings.Index(k, "_")] + "][" + k[strings.Index(k, "_")+1:len(k)] + "]"] = submission[k]
        } else {
            data["submission[" + k + "]"] = submission[k]   
        }
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formId, 10) + "/submissions", data, "POST")
}

//CreateFormSubmissions
//Submit data to this form using the API
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//submission (map[string]string): Submission data with question IDs.
//Returns posted submission ID and URL.
func (client jotformAPIClient) CreateFormSubmissions(formId int64, submission []byte) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formId, 10) + "/submissions", submission, "PUT")
}

//GetFormFiles
//List of files uploaded on a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns uploaded file information and URLs on a specific form.
func (client jotformAPIClient) GetFormFiles(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/files", "", "GET")
}

//GetFormWebhooks
//Get list of webhooks for a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns list of webhooks for a specific form.
func (client jotformAPIClient) GetFormWebhooks(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/webhooks", "", "GET")
}

//CreateFormWebhook
//Add a new webhook
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//webhookURL (string): Webhook URL is where form data will be posted when form is submitted.
//Returns list of webhooks for a specific form.
func (client jotformAPIClient) CreateFormWebhook(formId int64, webhookURL string) []byte {
    params := map[string]string {
        "webhookURL": webhookURL,
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formId, 10) + "/webhooks", params, "POST")
}

//Delete a specific webhook of a form.
//Add a new webhook
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//webhookID (int64): You can get webhook IDs when you call /form/{formID}/webhooks.
//Returns remaining webhook URLs of form.
func (client jotformAPIClient) DeleteFormWebhook(formID int64, webhookID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/webhooks/" + strconv.FormatInt(webhookID, 10), nil, "DELETE")
}

//GetSubmission
//Get submission data
//sid (int64): You can get submission IDs when you call /form/{id}/submissions.
//Returns information and answers of a specific submission.
func(client jotformAPIClient) GetSubmission(sid int64) []byte {
    return client.executeHttpRequest("user/submission/" + strconv.FormatInt(sid, 10), "","GET")
}

//GetReport
//Get report details
//reportID (int64): You can get a list of reports from /user/reports.
//Returns properties of a speceific report like fields and status.
func(client jotformAPIClient) GetReport(reportID int64) []byte {
    return client.executeHttpRequest("user/report/" + strconv.FormatInt(reportID, 10), "", "GET")
}

//GetFolder
//folderID (int64): You can get a list of folders from /user/folders.
//Returns a list of forms in a folder, and other details about the form such as folder color.
func (client jotformAPIClient) GetFolder(folderID string) []byte {
    return client.executeHttpRequest("folder/" + folderID, "", "GET")
}

//GetFormProperties
//Get a list of all properties on a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns form properties like width, expiration date, style etc.
func (client jotformAPIClient) GetFormProperties(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", "", "GET")
}

//GetFormReports
//Get all the reports of a form, such as excel, csv, grid, html, etc.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns list of all reports in a form, and other details about the reports such as title.
func (client jotformAPIClient) GetFormReports(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/reports", "", "GET")
}

//CreateReport
//Create new report of a form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//report (map[string]string): Report details. List type, title etc.
//Returns report details and URL.
func (client jotformAPIClient) CreateReport(formID int64, report map[string]string) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/reports", report, "POST")
}

//GetFormProperty
//Get a specific property of the form.]
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//propertyKey (string): You can get property keys when you call /form/{id}/properties.
//Returns given property key value.
func (client jotformAPIClient) GetFormProperty(formID int64, propertyKey string) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties/" + propertyKey, "", "POST")
}

//DeleteSubmission
//Delete a single submission
//sid (int64): You can get submission IDs when you call /form/{id}/submissions.
//Returns status of request.
func (client jotformAPIClient) DeleteSubmission(sid int64) []byte {
    return client.executeHttpRequest("submission/" + strconv.FormatInt(sid, 10), nil, "DELETE")
}

//EditSubmission
//Edit a single submission
//sid (int64): You can get submission IDs when you call /form/{id}/submissions.
//submission (map[string]string): New submission data with question IDs.
//Returns status of request.
func (client jotformAPIClient) EditSubmission(sid int64, submission map[string]string) []byte {
    data := make(map[string]string)

    for k, _ := range submission {
        if strings.Contains(k, "_") && k!= "created_at" {
            data["submission[" + k[0:strings.Index(k, "_")] + "][" + k[strings.Index(k, "_")+1:len(k)] + "]"] = submission[k]
        } else {
            data["submission[" + k + "]"] = submission[k]   
        }
    }

    return client.executeHttpRequest("submission/" + strconv.FormatInt(sid, 10), data, "POST")
}

//CloneForm
//Clone a single form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns status of request.
func (client jotformAPIClient) CloneForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/clone", nil, "POST")
}

//DeleteFormQuestion
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//qid (int): Identifier for each question on a form. You can get a list of question IDs from /form/{id}/questions.
//Returns status of request.
func (client jotformAPIClient) DeleteFormQuestion(formID int64, qid int) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), nil, "DELETE")
}

//CreateFormQuestion
//Add new question to specified form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//questionProperties (map[string]string): New question properties like type and text.
//Returns properties of new question.
func (client jotformAPIClient) CreateFormQuestion(formID int64, questionProperties map[string]string) []byte {
    question := make(map[string]string)

    for k, _ := range questionProperties {
        question["question[" + k + "]"] = questionProperties[k]
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", question, "POST")
}

//CreateFormQuestion
//Add new question to specified form.
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//questions ([]byte): New question properties like type and text.
//Returns properties of new question.
func (client jotformAPIClient) CreateFormQuestions(formID int64, questions []byte) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", questions, "PUT")
}

//EditFormQuestion
//Add or edit a single question properties
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//qid (int): Identifier for each question on a form. You can get a list of question IDs from /form/{id}/questions.
//questionProperties (map[string]string): New question properties like type and text.
//Returns edited property and type of question.
func (client jotformAPIClient) EditFormQuestion(formID int64, qid int, questionProperties map[string]string) []byte {
    question := make(map[string]string)

    for k, _ := range questionProperties {
        question["question[" + k + "]"] = questionProperties[k]
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), question, "POST")
}

//SetFormProperties
//Add or edit properties of a specific form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//formProperties (map[string]string): New properties like label width.
//Returns edited properties.
func (client jotformAPIClient) SetFormProperties(formID int64, formProperties map[string]string) []byte {
    properties := make(map[string]string)

    for k, _ := range formProperties {
        properties["properties[" + k + "]"] = formProperties[k]
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", properties, "POST")
}

//SetFormProperties
//Add or edit properties of a specific form
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//formProperties ([]byte): New properties like label width.
//Returns edited properties.
func (client jotformAPIClient) SetMultipleFormProperties(formID int64, formProperties []byte) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", formProperties, "PUT")
}

//CreateForm
//Create a new form
//form ([]byte): Questions, properties and emails of new form.
//Returns new form.
func (client jotformAPIClient) CreateForm(form map[string]interface{}) []byte{
    params := make(map[string]string)

    for formKey, formValue := range form {
            if formKey == "properties" {
                    properties := formValue

                    for properyKey, propertyValue := range properties.(map[string]string) {
                            params[formKey + "[" + properyKey + "]"] = propertyValue
                    }
            } else {
                    formItem := formValue

                    for formItemKey, formItemValue := range formItem.(map[string]interface{}) {
                            item := formItemValue

                            for itemKey, itemValue := range item.(map[string]string) {
                                    params[formKey + "[" + formItemKey + "][" + itemKey + "]"] = itemValue
                            }
                    }
            }
    }

    return client.executeHttpRequest("user/forms", params, "POST")
}

//Create new forms
//Create a new form
//form ([]byte): Questions, properties and emails of forms.
//Returns new forms.
func (client jotformAPIClient) CreateForms(form []byte) []byte {
    return client.executeHttpRequest("user/forms", form, "PUT")
}

//DeleteForm
//formID (int64): Form ID is the numbers you see on a form URL. You can get form IDs when you call /user/forms.
//Returns properties of deleted form.
func (client jotformAPIClient) DeleteForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10), nil, "DELETE")
}

//RegisterUser
//Register with username, password and email
//userDetails (map[string]string): Username, password and email to register a new user
//Returns new user's details
func (client jotformAPIClient) RegisterUser(userDetails map[string]string) []byte {
    return client.executeHttpRequest("user/register", userDetails, "POST")
}

//LoginUser
//Login user with given credentials
//credentials (map[string]string): Username, password, application name and access type of user
//Returns logged in user's settings and app key
func (client jotformAPIClient) LoginUser(credentials map[string]string) []byte {
    return client.executeHttpRequest("user/login", credentials, "POST");
}

//GetPlan
//Get details of a plan
//planName (string): Name of the requested plan. FREE, PREMIUM etc.
//Returns details of a plan
func (client jotformAPIClient) GetPlan(planName string) []byte {
    return client.executeHttpRequest("system/plan/" + planName, "", "GET")
}


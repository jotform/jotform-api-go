package jotform

import(
    "fmt"
    "strconv"
    "net/http"
    "net/url"   
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
)

const baseURL = "http://api.jotform.com"
const apiVersion = "v1"

type jotformAPIClient struct{
    apiKey string
}

func NewJotFormAPIClient(apiKey string) *jotformAPIClient {
    client := &jotformAPIClient{apiKey}

    return client
}

func (client jotformAPIClient) checkClient() {
    if client.apiKey == "" {
        fmt.Printf("ApiKey is requied!\n")
        os.Exit(1)
    }
}

func (client jotformAPIClient) executeHttpRequest(requestPath string, params interface{}, method string) []byte {
    client.checkClient()

    var path = baseURL + "/" + apiVersion + "/" + requestPath

    var response *http.Response
    var request *http.Request
    var err error

    if method == "GET" {
        path = path + "?" + params.(string) + "apiKey=" + client.apiKey

        response, err = http.Get(path)

    } else if method == "POST" {
        path = path + "?" + "apiKey=" + client.apiKey

        data := params.(map[string]string)

        values := make(url.Values)

        for k, _ := range data {
            values.Set(k, data[k])
        }

        response, err = http.PostForm(path, values)

    } else  if method == "DELETE" {
        path = path + "?" + "apiKey=" + client.apiKey

        request, err = http.NewRequest("DELETE", path, nil)

        response, err = http.DefaultClient.Do(request)
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

        var f interface{}

        json.Unmarshal(contents, &f)

        result := f.(map[string]interface{})["content"]

        content, err := json.Marshal(result)

        return content
    }

    return nil
}

func createConditions (offset string, limit string, filter map[string]string, orderBy string) string {
    var params = ""

    if offset != "" {
        params = "offset=" + offset + "&"
    }

    if limit!= "" {
        params = params + "limit=" + limit + "&"
    }

    if filter != nil {
        var value = "{"
        var count = 0

        for k, _ := range filter {
            count++

            value = value + "\"" + k + "\":\"" + filter[k] + "\""

            if count < len(filter) {
                value = value + ","
            }
        }
        value = value + "}&"

        params = params + "filter=" + value
    }

    if orderBy != "" {
        params = params + "order_by=" + orderBy + "&"
    }

    return params
}

func (client jotformAPIClient) GetUser() []byte {
    return client.executeHttpRequest("user", "", "GET")
}

func (client jotformAPIClient) GetUsage() []byte {
    return client.executeHttpRequest("user/usage", "", "GET")
}

func (client jotformAPIClient) GetForms(offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("user/forms", params, "GET")
}

func (client jotformAPIClient) GetSubmissions(offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("user/submissions", params, "GET")
}

func (client jotformAPIClient) GetSubusers() []byte {
    return client.executeHttpRequest("user/subusers", "", "GET")
}

func (client jotformAPIClient) GetFolders() []byte {
    return client.executeHttpRequest("user/folders", "", "GET")
}

func (client jotformAPIClient) GetReports() []byte {
    return client.executeHttpRequest("user/reports", "", "GET")
}

func (client jotformAPIClient) GetSettings() []byte {
    return client.executeHttpRequest("user/settings", "", "GET")
}

func (client jotformAPIClient) GetHistory() []byte {
    return client.executeHttpRequest("user/history", "", "GET")
}

func (client jotformAPIClient) GetForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10), "", "GET")
}

func (client jotformAPIClient) GetFormQuestions(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/questions", "", "GET")
}

func (client jotformAPIClient) GetFormQuestion(formID int64, qid int) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), "", "GET")
}

func (client jotformAPIClient) GetFormSubmissions(formID int64, offset string, limit string, filter map[string]string, orderBy string) []byte {
    var params = createConditions(offset, limit, filter, orderBy)

    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/submissions", params, "GET")
}

func (client jotformAPIClient) CreateFormSubmissions(formId int64, submission map[string]string) []byte {
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

func (client jotformAPIClient) GetFormFiles(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/files", "", "GET")
}

func (client jotformAPIClient) GetFormWebhooks(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/webhooks", "", "GET")
}

func (client jotformAPIClient) CreateFormWebhook(formId int64, webhookURL string) []byte {
    params := map[string]string {
        "webhookURL": webhookURL,
    }

    return client.executeHttpRequest("form/" + strconv.FormatInt(formId, 10) + "/webhooks", params, "POST")
}

func(client jotformAPIClient) GetSubmission(sid int64) []byte {
    return client.executeHttpRequest("user/submission/" + strconv.FormatInt(sid, 10), "","GET")
}

func(client jotformAPIClient) GetReport(reportID int64) []byte {
    return client.executeHttpRequest("user/report/" + strconv.FormatInt(reportID, 10), "", "GET")
}

func (client jotformAPIClient) GetFolder(folderID int64) []byte {
    return client.executeHttpRequest("user/folder/" + strconv.FormatInt(folderID, 10), "", "GET")
}

func (client jotformAPIClient) GetFormProperties(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties", "", "GET")
}

func (client jotformAPIClient) GetFormProperty(formID int64, propertyKey string) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/properties/" + propertyKey, "", "POST")
}

func (client jotformAPIClient) DeleteSubmission(sid int64) []byte {
    return client.executeHttpRequest("submission/" + strconv.FormatInt(sid, 10), nil, "DELETE")
}

func (client jotformAPIClient) EditSubmission(sid int64, submission map[string]string) []byte {
    data := make(map[string]string)

    for k, _ := range submission {
        if strings.Contains(k, "_") {
            data["submission[" + k[0:strings.Index(k, "_")] + "][" + k[strings.Index(k, "_")+1:len(k)] + "]"] = submission[k]
        } else {
            data["submission[" + k + "]"] = submission[k]   
        }
    }

    return client.executeHttpRequest("submission/" + strconv.FormatInt(sid, 10), data, "POST")
}

func (client jotformAPIClient) CloneForm(formID int64) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/clone", nil, "POST")
}

func (client jotformAPIClient) DeleteFormQuestion(formID int64, qid int) []byte {
    return client.executeHttpRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid), nil, "DELETE")
}


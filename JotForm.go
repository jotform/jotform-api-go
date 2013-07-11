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

type JotformAPIClient struct{
    Username string
    ApiKey string
}

func (client JotformAPIClient) checkClient() {
    if client.ApiKey == "" {
        fmt.Printf("ApiKey is requied!\n")
        os.Exit(1)
    }
}

func (client JotformAPIClient) getRequest(requestPath string) []byte {
    client.checkClient()

    var path = baseURL + "/" + apiVersion + "/" + requestPath + "?apiKey=" + client.ApiKey

    response, err := http.Get(path)
    
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

func (client JotformAPIClient) postRequest(requestPath string, params map[string]string) [] byte {
    client.checkClient()

    var path = baseURL + "/" + apiVersion + "/" + requestPath + "?apiKey=" + client.ApiKey


    values := make(url.Values)

    for k, _ := range params {
        values.Set(k, params[k])
    }

    response, err := http.PostForm(path, values)

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

func (client JotformAPIClient) GetUser() []byte {
    return client.getRequest("user")
}

func (client JotformAPIClient) GetUsage() []byte {
    return client.getRequest("user/usage")
}

func (client JotformAPIClient) GetForms() []byte {
    return client.getRequest("user/forms")
}

func (client JotformAPIClient) GetSubmissions() []byte {
    return client.getRequest("user/submissions")
}

func (client JotformAPIClient) GetSubusers() []byte {
    return client.getRequest("user/subusers")
}

func (client JotformAPIClient) GetFolders() []byte {
    return client.getRequest("user/folders")
}

func (client JotformAPIClient) GetReports() []byte {
    return client.getRequest("user/reports")
}

func (client JotformAPIClient) GetSettings() []byte {
    return client.getRequest("user/settings")
}

func (client JotformAPIClient) GetHistory() []byte {
    return client.getRequest("user/history")
}

func (client JotformAPIClient) GetForm(formID int64) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10))
}

func (client JotformAPIClient) GetFormQuestions(formID int64) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/questions")
}

func (client JotformAPIClient) GetFormQuestion(formID int64, qid int) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/question/" + strconv.Itoa(qid))
}

func (client JotformAPIClient) GetFormSubmissions(formID int64) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/submissions")
}

func (client JotformAPIClient) CreateFormSubmissions(formId int64, submission map[string]string) []byte {
    data := make(map[string]string)

    for k, _ := range submission {
        if strings.Contains(k, "_") {
            data["submission[" + k[0:strings.Index(k, "_")] + "][" + k[strings.Index(k, "_")+1:len(k)] + "]"] = submission[k]
        } else {
            data["submission[" + k + "]"] = submission[k]   
        }
    }

    return client.postRequest("form/" + strconv.FormatInt(formId, 10) + "/submissions", data)
}

func (client JotformAPIClient) GetFormFiles(formID int64) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/files")
}

func (client JotformAPIClient) GetFormWebhooks(formID int64) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/webhooks")
}

func (client JotformAPIClient) CreateFormWebhook(formId int64, webhookURL string) []byte {
    params := map[string]string {
        "webhookURL": webhookURL,
    }

    return client.postRequest("form/" + strconv.FormatInt(formId, 10) + "/webhooks", params)
}

func(client JotformAPIClient) GetSubmission(sid int64) []byte {
    return client.getRequest("user/submission/" + strconv.FormatInt(sid, 10))
}

func(client JotformAPIClient) GetReport(reportID int64) []byte {
    return client.getRequest("user/report/" + strconv.FormatInt(reportID, 10))
}

func (client JotformAPIClient) GetFolder(folderID int64) []byte {
    return client.getRequest("user/folder/" + strconv.FormatInt(folderID, 10))
}

func (client JotformAPIClient) GetFormProperties(formID int64) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/properties")
}

func (client JotformAPIClient) GetFormProperty(formID int64, propertyKey string) []byte {
    return client.getRequest("form/" + strconv.FormatInt(formID, 10) + "/properties/" + propertyKey)
}
jotform-api-go
==============
[JotForm API](http://api.jotform.com/docs/) - GO Client

### Installation

Install via git clone:

        $ git clone git://github.com/jotform/jotform-api-go.git
        $ cd jotform-api-go

### Documentation

You can find the docs for the API of this client at [http://api.jotform.com/docs/](http://api.jotform.com/docs)

### Authentication

JotForm API requires API key for all user related calls. You can create your API Keys at  [API section](http://www.jotform.com/myaccount/api) of My Account page.

### Examples

Get latest 100 submissions ordered by creation date

```go
package main

import (
        "jotform-api-go"
        "fmt"
)

func main() {

    jotformAPI := jotform.NewJotFormAPIClient("YOUR API KEY")

    submissions := jotformAPI.GetSubmissions("", "100", nil, "created_at")
    fmt.Println(string(submissions))
}
``` 
    
Submission and form filter examples

```go
package main

import (
        "jotform-api-go"
        "fmt"
)

func main() {

    jotformAPI := jotform.NewJotFormAPIClient("YOUR API KEY")


    submissionFilter := map[string]string {
            "id:gt": "FORM ID",
            "created_at:gt": "DATE",
    }

    submissions := jotformAPI.GetSubmissions("", "", submissionFilter, "")
    fmt.Println(string(submissions))

    formFilter := map[string]string {
            "id:gt": "FORM ID",       
    }

    forms := jotformAPI.GetForm("", "", formFilter, "")
    fmt.Println(string(forms))
}
``` 

First the _jotform_ package is imported from the _jotform-api-go/JotForm.go_ file. This package provides access to JotForm's API. You have to create an API client instance with your API key. 
In any case of exception (wrong authentication etc.), you can catch it or let it fail with fatal error.
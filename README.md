jotform-api-go
==============
[JotForm API](https://api.jotform.com/docs/) - GO Client

**Strongly consider using [v2 of the jotform-api-go client!](https://github.com/jotform/jotform-api-go/tree/master/v2)!**

v1 works for user scripts, but is unsafe for use in long-running servers.
It also silently hides some errors, and is generally less useable than v2.

### Installation

Install via git clone:

        $ git clone git://github.com/jotform/jotform-api-go.git
        $ cd jotform-api-go

### Documentation

You can find the docs for the API of this client at [https://api.jotform.com/docs/](https://api.jotform.com/docs)

### Authentication

JotForm API requires API key for all user related calls. You can create your API Keys at  [API section](https://www.jotform.com/myaccount/api) of My Account page.

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

    jotformAPI := jotform.NewJotFormAPIClient(
        "YOUR API KEY",
        "json", // or "xml", depending on how you want results to be formatted
        false,  // when true, this prints debugging information on each request
    )


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
In case of an exception (wrong authentication etc.), you can catch it or let it fail with a fatal error.

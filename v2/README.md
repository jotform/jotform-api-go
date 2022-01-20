jotform-api-go
==============
[JotForm API](http://api.jotform.com/docs/) - GO Client

**The v2 client should not yet be considered stable.**

In particular, changes are expected to all non-download interfaces
(to return structs with the expected fields,
rather than relying on the consumer to know how to parse eg. JotForm's json)
and `NewJotFormClient()`, reflecting the above direction.

Error reporting can also be made more robust.

### Installation

This can be installed the same as any other go module.
Using [go modules](https://go.dev/blog/using-go-modules),
this is as simple as `go get github.com/jotform/jotform-api-go/v2`.

This will let you import the module into a go package, like so:

```go
import (
        "fmt"
        jotform "github.com/jotform/jotform-api-go/v2"
)
```

### Documentation

You can find the docs for the API of this client at [http://api.jotform.com/docs/](http://api.jotform.com/docs)

### Authentication

JotForm API requires API key for all user related calls. You can create your API Keys at  [API section](http://www.jotform.com/myaccount/api) of My Account page.

### Examples

Get submission ID 1234567.

```go
package main

import (
        "fmt"
        jotform "github.com/jotform/jotform-api-go/v2"
)

func main() {

    jotformAPI := jotform.NewJotFormAPIClient(
        "YOUR API KEY",
        "json",
        false,
    )

    // NOTE the structure of submission in this response
    // is likely to change in the future
    submission, err := jotformAPI.GetSubmission(int64(1234567))
    if err != nil {
        ...
    }

    fmt.Println(string(submission))
}
```

### Testing

You can run the tests for v2 like so:

```
$ cd v2 && go test ./...
```

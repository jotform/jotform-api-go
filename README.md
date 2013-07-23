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


First the jotform package is imported from the jotform-api-go/JotForm.go file. This package provides access to JotForm's API. You have to create an API client instance with your API key. 
In any case of exception (wrong authentication etc.), you can catch it or let it fail with fatal error.
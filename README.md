# web-search

[![Codeship Status for ihcsim/web-search](https://app.codeship.com/projects/078c5ba0-b663-0134-52cc-7ea8c0f9c13a/status?branch=master)](https://app.codeship.com/projects/194346)

A golang web-search server that submits search requests to the Google Custom Search API.

This project requires Golang 1.7. It explores the following Go features:
* The [context](https://golang.org/pkg/context/) package to propagate request-scoped data to the remote API. web-search stores the requests originating IP address and request timeout in context.
  * Refer this [post](https://blog.golang.org/context) for more information on the `context` package.
  * For more information on using context to cancel requests, refer:
    * https://golang.org/pkg/net/http/#Request
    * https://golang.org/pkg/net/http/#Request.WithContext
    * https://golang.org/pkg/net/http/#Request.Context
* The [github.com/pkg/errors](https://godoc.org/github.com/pkg/errors) package to wrap errors in stacktraces. Refer this [post](http://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully) for more information on annotating errors.
* The [httptest.ResponseRecorder](https://golang.org/pkg/net/http/httptest/#ResponseRecorder) struct to evaluate HTTP responses code and body in tests.
* The [Golang 1.7 subtests](https://golang.org/doc/go1.7#testing) to group tests into hierarchical structure.

To start web-search at 127.0.0.1:7000 with a request timeout of 2s:
```
$ go get github.com/ihcsim/web-search
$ SERVER_HOSTNAME=127.0.0.1:7000 REQUEST_TIMEOUT=2s go run cmd/server/main.go
```

To use curl to submit a search query:
```
$ curl -v 127.0.0.1:7000/search?q=golang
```

Supported environmental variables:

Variables | Description | Default
--------- | ----------- | -------
`SERVER_HOSTNAME` | The hostname and port number that web-search listens on | None
`REQUEST_TIMEOUT` | The timeout duration of a search request. | None

# LICENSE
Refer the [LICENSE](LICENSE) file.

package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ihcsim/web-search"
	"github.com/pkg/errors"
)

var (
	serverAddr string
	timeout    time.Duration
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if serverAddr = os.Getenv("SERVER_HOSTNAME"); serverAddr == "" {
		log.Fatal("Can't start server. Please specify the server's hostname.")
	}

	if t, err := time.ParseDuration(os.Getenv("REQUEST_TIMEOUT")); err != nil {
		log.Print("Starting server without request timeout. Error from package ", err)
	} else {
		timeout = t
	}

	http.HandleFunc("/search", searchHandler)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		handleErr(w, errors.Wrap(err, "Can't determine user IP"), http.StatusInternalServerError)
		return
	}

	ctx, cancel, err := searchContext(ipAddr)
	if err != nil {
		handleErr(w, errors.Wrap(err, "Can't create search context"), http.StatusInternalServerError)
		return
	}
	defer cancel()

	query := r.FormValue("q")
	if query == "" {
		handleErr(w, fmt.Errorf("Can't search. Please provide search query"), http.StatusBadRequest)
		return
	}

	start := time.Now()
	results, err := search.Google(ctx, query)
	elapsed := time.Since(start)
	if err != nil {
		handleErr(w, errors.Wrap(err, "Google Search API error."), http.StatusInternalServerError)
		return
	}

	if err := resultsTemplate.Execute(w, struct {
		Results          search.Results
		Timeout, Elapsed time.Duration
	}{
		Results: results,
		Timeout: timeout,
		Elapsed: elapsed,
	}); err != nil {
		log.Print(err)
		return
	}
}

func searchContext(ipAddr string) (ctx context.Context, cancel context.CancelFunc, err error) {
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	sourceIP := search.NewSourceIP(ipAddr)
	ctx, err = sourceIP.NewContext(ctx)
	return
}

func handleErr(w http.ResponseWriter, err error, code int) {
	log.Println(err.Error())
	http.Error(w, err.Error(), code)
}

var resultsTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
<ol>
{{range .Results}}
<li>{{.Title}} - <a href="{{.URL}}">{{.URL}}</a></li>
{{end}}
</ol>
<p>{{len .Results}} results in {{.Elapsed}}; timeout {{.Timeout}}</p>
</body>
</html>
`))

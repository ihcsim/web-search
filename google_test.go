package search

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

var searchDB = []struct {
	title string
	url   string
}{
	{title: "The Go Programming Language", url: "https://golang.org"},
	{title: "The Go Playground", url: "https://play.golang.org"},
	{title: "golang - GitHub", url: "https://github.com/golang"},
	{title: "r/golang - Reddit", url: "https://www.reddit.com/r/golang"},
	{title: "Go (programming language) - Wikipedia", url: "https://en.wikipedia.org/wiki/Go_(programming_language)"},
}

func TestGoogle(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	svr := httptest.NewServer(http.HandlerFunc(handleSearchRequest))
	apiURL = svr.URL

	query := "golang"
	ctx := context.WithValue(context.Background(), KeyIPAddr, SourceIP(net.ParseIP("10.0.0.1")))
	results, err := Google(ctx, query)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	for index, result := range results {
		if result.Title != searchDB[index].title {
			t.Errorf("Mismatch title. Expected %q, but got %q", searchDB[index].title, result.Title)
		}

		if result.URL != searchDB[index].url {
			t.Errorf("Mismatch URL. Expected %q, but got %q", searchDB[index].url, result.URL)
		}
	}
}

func handleSearchRequest(w http.ResponseWriter, r *http.Request) {
	result := bytes.Buffer{}
	for index, tuple := range searchDB {
		jsonObj := `{"titleNoFormatting":"%s", "url":"%s"}`
		if index < len(searchDB)-1 {
			jsonObj = jsonObj + ","
		}

		if _, err := result.Write([]byte(fmt.Sprintf(jsonObj, tuple.title, tuple.url))); err != nil {
			continue
		}
	}

	var serverResponse = fmt.Sprintf(`{"responseData": {"results": [%s]}}`, result.String())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(serverResponse))
}

func TestSearchRequest(t *testing.T) {
	sourceIP := NewSourceIP("172.0.0.1")
	ctx, err := sourceIP.NewContext(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	query := "golang"
	r, err := newSearchRequest(ctx, query)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	if r.Method != apiMethod {
		t.Errorf("Mismatched HTTP method. Expected %q, but got %q", apiMethod, r.Method)
	}

	actualURL := r.URL.Scheme + "://" + path.Join(r.URL.Host, r.URL.Path)
	if actualURL != apiURL {
		t.Errorf("Mismatched API URL. Expected %q, but got %q", apiURL, actualURL)
	}

	actualVersion := r.URL.Query().Get("v")
	if actualVersion != apiVersion {
		t.Errorf("Mismatched API version. Expected %q, but got %q", apiVersion, actualVersion)
	}

	actualQuery := r.URL.Query().Get("q")
	if actualQuery != query {
		t.Errorf("Mismatched query. Expected %q, but got %q", query, actualQuery)
	}

	actualIP := r.URL.Query().Get("userip")
	if actualIP != string(sourceIP) {
		t.Errorf("Mismatched source IP. Expected %q, but got %q", actualIP, actualIP)
	}
}

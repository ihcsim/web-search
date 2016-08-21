package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

var (
	apiMethod  = "GET"
	apiURL     = "https://ajax.googleapis.com/ajax/services/search/web"
	apiVersion = "1.0"
)

// Results contains an ordered list of search results.
type Results []Result

// Result contains the title and URL of a search result.
type Result struct {
	Title, URL string
}

// Google uses Google Custom Search to perform a search based on query.
// If ctx isn't needed, pass in context.TODO.
func Google(ctx context.Context, query string) (Results, error) {
	req, err := newSearchRequest(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create HTTP request")
	}

	// the request is closed when the context expired
	// Refer https://golang.org/pkg/net/http/#Request.Context
	req = req.WithContext(ctx)

	results := []Result{}
	err = httpDo(req, func(resp *http.Response, err error) error {
		if err != nil {
			return errors.Wrap(err, "Errors from Google Custom Search API")
		}
		defer resp.Body.Close()

		// Parse the JSON search result.
		// https://developers.google.com/web-search/docs/#fonje
		var data struct {
			ResponseData struct {
				Results []struct {
					TitleNoFormatting string
					URL               string
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.Wrap(err, "Failed to decode JSON response")
		}

		for _, res := range data.ResponseData.Results {
			results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
		}
		return nil
	})

	// httpDo waits for the closure we provided to return, so it's safe to
	// read results here.
	return results, err
}

func newSearchRequest(ctx context.Context, query string) (*http.Request, error) {
	req, err := http.NewRequest(apiMethod, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Google Search API failure")
	}

	q := req.URL.Query()
	q.Set("q", query)
	q.Set("v", apiVersion)

	var s SourceIP
	if err := s.FromContext(ctx); err != nil {
		return nil, errors.Wrap(err, "Can't set up search context")
	}
	q.Set("userip", fmt.Sprintf("%s", s))
	req.URL.RawQuery = q.Encode()

	return req, nil
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	search "github.com/ihcsim/web-search"
)

func TestSearchContext(t *testing.T) {
	if err := os.Setenv("SERVER_HOSTNAME", "127.0.0.1"); err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer func() {
		if err := os.Unsetenv("SERVER_HOSTNAME"); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	}()

	var tests = []struct {
		ipAddr   string
		expected search.SourceIP
	}{
		{ipAddr: "172.0.0.1", expected: search.SourceIP(net.ParseIP("172.0.0.1"))},
		{ipAddr: "10.0.0.1", expected: search.SourceIP(net.ParseIP("10.0.0.1"))},
		{ipAddr: "192.168.0.1", expected: search.SourceIP(net.ParseIP("192.168.0.1"))},
		{ipAddr: "2001:db8::68", expected: search.SourceIP(net.ParseIP("2001:db8::68"))},
		{ipAddr: "2001:db8:85a3:0:0:8a2e:370:7334", expected: search.SourceIP(net.ParseIP("2001:db8:85a3:0:0:8a2e:370:7334"))},
		{ipAddr: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", expected: search.SourceIP(net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))},
	}

	for _, test := range tests {
		ctx, cancel, err := searchContext(test.ipAddr)
		if err != nil {
			t.Fatal("Unexpected error occurred: ", err)
		}
		defer cancel()

		actual, ok := ctx.Value(search.KeyIPAddr).(search.SourceIP)
		if !ok {
			t.Fatal("Unexpected type assertion failure")
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Mismatch IP. Expected %v, but got %v", test.expected, actual)
		}
	}
}

func TestHandleErr(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	var tests = []struct {
		w    *httptest.ResponseRecorder
		code int
		err  string
	}{
		{w: httptest.NewRecorder(), code: http.StatusBadRequest, err: "Bad Request"},
		{w: httptest.NewRecorder(), code: http.StatusForbidden, err: "Forbidden"},
		{w: httptest.NewRecorder(), code: http.StatusInternalServerError, err: "Internal Server Error"},
	}

	for _, test := range tests {
		handleErr(test.w, fmt.Errorf(test.err), test.code)

		if test.w.Code != test.code {
			t.Errorf("Mismatch status. Expected %v, but got %v", test.code, test.w.Code)
		}

		actualErr := strings.TrimSuffix(test.w.Body.String(), "\n")
		if actualErr != test.err {
			t.Errorf("Mismatch error body. Expected %q, but got %q", test.err, actualErr)
		}
	}
}

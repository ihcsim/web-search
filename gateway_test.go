package search

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const errMsg = "An error has occurred"

func TestHTTPDo(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("When response handler returns no errors", func(t *testing.T) {
		if err := httpDo(&http.Request{}, handleNoErrorResponse); err != nil {
			t.Error("Unexpected error: ", err)
		}
	})

	t.Run("When response handler returns errors", func(t *testing.T) {
		err := httpDo(&http.Request{}, handleErrorResponse)
		if err == nil {
			t.Errorf("Expected error with message %q didn't occur", errMsg)
		}

		if err.Error() != errMsg {
			t.Errorf("Mismatched error message. Expected %q, but got %q", errMsg, err)
		}
	})

	t.Run("When context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		r := (&http.Request{}).WithContext(ctx)

		go func() {
			// this gorountine cancels the context after 1ms because
			// handleNoResponse() will never return
			<-time.After(time.Millisecond)
			cancel()
		}()

		err := httpDo(r, handleNoResponse)
		if err == nil {
			t.Errorf("Expected error with message %q didn't occur", context.Canceled)
		}

		if err != context.Canceled {
			t.Errorf("Mismatched error message. Expected %q, but got %q", context.Canceled, err)
		}
	})

	t.Run("When context timed out", func(t *testing.T) {
		timeout := time.Microsecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		// set the request context with the specified timeout
		r := (&http.Request{}).WithContext(ctx)

		err := httpDo(r, handleNoResponse)
		if err == nil {
			t.Errorf("Expected error with message %q didn't occur", context.DeadlineExceeded)
		}

		if err != context.DeadlineExceeded {
			t.Errorf("Mismatched error message. Expected %s, but got %s", context.DeadlineExceeded, err)
		}
	})
}

func handleNoErrorResponse(res *http.Response, err error) error {
	return nil
}

func handleErrorResponse(res *http.Response, err error) error {
	return fmt.Errorf("%s", errMsg)
}

func handleNoResponse(res *http.Response, err error) error {
	for {
		// runs infinitely to trigger context timeout
	}
	return nil
}

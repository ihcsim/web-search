package search

import (
	"context"
	"log"
	"net/http"
)

const msgResponseReceived = "Received response from %s"

func httpDo(req *http.Request, handleResponse func(*http.Response, error) error) error {
	client := &http.Client{}
	c := make(chan error, 1)

	go func() {
		c <- handleResponse(client.Do(req))
	}()

	select {
	case <-req.Context().Done():
		log.Println(context.Canceled.Error())
		return req.Context().Err()
	case err := <-c:
		log.Printf(msgResponseReceived, req.URL)
		return err
	}
}

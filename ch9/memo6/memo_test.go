// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
package memo6

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

func httpGetBody(url string, done <-chan struct{}) (interface{}, error) {
	mainContext := context.Background()
	newContext, cancelFunction := context.WithCancel(mainContext)
	go func() {
		<-done
		cancelFunction()
	}()
	request, err := http.NewRequestWithContext(newContext, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func TestCancellation(t *testing.T) {
	m := New(httpGetBody)
	done := make(chan struct{})
	close(done)
	url := "https://golang.org"
	_, err := m.Get(url, done)
	if err == nil {
		t.Error(err)
	}
}

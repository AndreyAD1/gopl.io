// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
package memo6

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
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
		t.Fatal("expect error, got nil")
	}
	expectedError := `Get "https://golang.org": context canceled`
	if err.Error() != expectedError {
		t.Errorf(
			"expected error message: %s, got: %s",
			expectedError,
			err.Error(),
		)
	}
}

type testCase struct {
	url          string
	cancellation bool
	expectError  bool
}

func incomingURLs() <-chan testCase {
	ch := make(chan testCase)
	go func() {
		for _, url := range []testCase{
			{"https://golang.org", false, false},
			{"https://godoc.org", true, true},
			{"https://play.golang.org", true, true},
			{"http://gopl.io", false, false},
			{"https://golang.org", false, false},
			{"https://godoc.org", true, true},
			{"https://play.golang.org", false, false},
			{"http://gopl.io", true, false},
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

func TestSequentialWithCancellation(t *testing.T) {
	m := New(httpGetBody)
	for testCase := range incomingURLs() {
		t.Run(testCase.url, func(t *testing.T) {
			start := time.Now()
			done := make(chan struct{})
			if testCase.cancellation {
				close(done)
			}
			value, err := m.Get(testCase.url, done)
			if testCase.expectError && err == nil {
				t.Fatalf("expect error, got nil for URL %s", testCase.url)
			}
			if err != nil {
				expectedError := fmt.Sprintf(
					`Get "%s": context canceled`, 
					testCase.url,
				)
				if err.Error() != expectedError {
					t.Fatalf(
						"expected error message: %s, got: %s",
						err.Error(),
						expectedError,
					)
				}
			}
			if !testCase.cancellation && value == nil {
				t.Fatalf("empty response body for URL %s", testCase.url)
			}
			if value != nil {
				fmt.Printf(
					"%s, %s, %d bytes\n",
					testCase.url,
					time.Since(start),
					len(value.([]byte)),
				)
			}
		})
	}
}

func TestConcurrentWithCancellation(t *testing.T) {
	m := New(httpGetBody)
	for testCase := range incomingURLs() {
		t.Run(testCase.url, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			done := make(chan struct{})
			if testCase.cancellation {
				close(done)
			}
			value, err := m.Get(testCase.url, done)
			if testCase.cancellation && err == nil {
				t.Fatalf("expect error, got nil for URL %s", testCase.url)
			}
			if testCase.cancellation && err != nil {
				expectedError := fmt.Sprintf(
					`Get "%s": context canceled`, 
					testCase.url,
				)
				if err.Error() != expectedError {
					t.Fatalf(
						"expected error message: %s, got: %s",
						err.Error(),
						expectedError,
					)
				}
			}
			if !testCase.cancellation && value == nil {
				t.Fatalf("empty response body for URL %s", testCase.url)
			}
			if value != nil {
				fmt.Printf(
					"%s, %s, %d bytes\n",
					testCase.url,
					time.Since(start),
					len(value.([]byte)),
				)
			}
		})
	}
}

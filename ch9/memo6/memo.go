package memo6

import (
	"fmt"
	"sync"
)

// Func is the type of the function to memoize.
type Func func(string, <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

//!+
type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string, done <-chan struct{}) (value interface{}, err error) {
	resultChannel := make(chan *entry)
	go memo.get(key, done, resultChannel)
	select {
	case result := <-resultChannel:
		return result.res.value, result.res.err
	case <-done:
		return nil, fmt.Errorf("receive a cancel signal")
	}
}

func (memo *Memo) get(key string, done <-chan struct{}, resultChannel chan<- *entry) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()
		value, err := memo.f(key, done)
		select {
		case <-done:
			memo.mu.Lock()
			delete(memo.cache, key)
			memo.mu.Unlock()
		default:
			e.res.value, e.res.err = value, err
		}

		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	resultChannel<- e
}

//!-


package memo7

type Func func(key string, done <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type request struct {
	key      string
	done     <-chan struct{}
	response chan<- result
}

type Memo struct{ requests chan request }

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, done, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

// type cancelEvent struct {
// 	isCanceled  bool
// 	canceledKey string
// }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	cancellationChannel := make(chan string)
	for req := range memo.requests {
		select {
		case cancelledKey := <-cancellationChannel:
			delete(cache, cancelledKey)
		default:
		}
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, req.done, cancellationChannel)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(
	f Func, 
	key string, 
	done <-chan struct{}, 
	cancellationChannel chan<- string,
) {
	e.res.value, e.res.err = f(key, done)
	close(e.ready)
	select {
	case <-done:
		cancellationChannel<- key
	default:
	}
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}

package queue

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RequestResult holds the result of a queued request.
type RequestResult struct {
	Body []byte
	Err  error
}

// RequestFunc is a function that performs an HTTP request using the provided client,
// and sends its result to the given result channel.
type RequestFunc func(client *http.Client, done chan<- RequestResult)

type queuedRequest struct {
	job  RequestFunc
	done chan RequestResult
}

type RequestQueue struct {
	client  *http.Client
	limiter *rate.Limiter
	queue   chan queuedRequest
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	closed  sync.Once
}

// NewRequestQueue creates a new rate-limited request queue.
// `rps` = requests per second, `burst` = burst capacity
func NewRequestQueue(rps int, burst int) *RequestQueue {
	if rps < 1 {
		rps = 1
	}
	if burst < 1 {
		burst = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	q := &RequestQueue{
		client:  &http.Client{Timeout: 10 * time.Second},
		limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), burst),
		queue:   make(chan queuedRequest, 100),
		ctx:     ctx,
		cancel:  cancel,
	}
	q.start()
	return q
}

// Enqueue adds a new request to the queue and returns a result channel.
func (q *RequestQueue) Enqueue(job RequestFunc) <-chan RequestResult {
	resultChan := make(chan RequestResult, 1)
	defer func() {
		if recover() != nil {
			resultChan <- RequestResult{Err: context.Canceled}
		}
	}()

	if q.ctx.Err() != nil {
		resultChan <- RequestResult{Err: context.Canceled}
		return resultChan
	}

	q.queue <- queuedRequest{
		job:  job,
		done: resultChan,
	}

	return resultChan
}

// start begins the worker loop.
func (q *RequestQueue) start() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Println("[queue] worker goroutine panicked:", r)
			}
		}()
		for req := range q.queue {
			if err := q.limiter.Wait(q.ctx); err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					req.done <- RequestResult{Err: context.Canceled}
					for pending := range q.queue {
						pending.done <- RequestResult{Err: context.Canceled}
					}
					return
				}
				log.Println("[queue] limiter wait failed:", err)
				req.done <- RequestResult{Err: err}
				continue
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						req.done <- RequestResult{Err: errors.New("request job panicked")}
					}
				}()
				req.job(q.client, req.done)
			}()
		}
	}()
}

// Shutdown stops the queue and waits for the worker to finish.
func (q *RequestQueue) Shutdown() {
	q.closed.Do(func() {
		q.cancel()
		close(q.queue)
		q.wg.Wait()
	})
}

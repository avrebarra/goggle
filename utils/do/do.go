package do

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

func Parallel(ctx context.Context, concurrency int, fns []func() error) (errs []error) {
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(concurrency))
	errChan := make(chan error, len(fns))

	for _, fn := range fns {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := sem.Acquire(ctx, 1); err != nil {
				errChan <- err
				return
			}
			defer sem.Release(1)
			errChan <- fn()
		}(fn)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	errs = make([]error, len(fns))
	for i := 0; i < len(fns); i++ {
		errs[i] = <-errChan
	}

	return errs
}

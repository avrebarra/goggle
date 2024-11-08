package do

import (
	"context"
	"sync"
)

func Parallel(ctx context.Context, fns []func() error) (errs []error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(fns))
	resultOrder := make(chan int, len(fns))

	for i, fn := range fns {
		wg.Add(1)
		go func(i int, fn func() error) {
			defer wg.Done()
			errChan <- fn()
			resultOrder <- i
		}(i, fn)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(resultOrder)
	}()

	errs = make([]error, len(fns))
	for i := range resultOrder {
		err := <-errChan
		errs[i] = err
	}

	return errs
}

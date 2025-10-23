package procedimientos

import (
	"fmt"
	"strings"
	"sync"
)

// runAdaptive attempts to run the provided task functions in parallel. If any
// of them return an error, runAdaptive will fall back to running them
// sequentially (one-by-one) and aggregate errors. It returns an aggregated
// error (or nil).
func runAdaptive(tasks []func() error) error {
	if len(tasks) == 0 {
		return nil
	}

	// try parallel first using WaitGroup + error channel
	var wg sync.WaitGroup
	errs := make(chan error, len(tasks))
	for _, t := range tasks {
		t := t
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- t()
		}()
	}
	wg.Wait()
	close(errs)

	var anyErr bool
	for e := range errs {
		if e != nil {
			anyErr = true
			break
		}
	}
	if !anyErr {
		return nil
	}

	// fallback: run sequentially and collect errors
	var parts []string
	for i, t := range tasks {
		if err := t(); err != nil {
			parts = append(parts, fmt.Sprintf("task[%d]: %v", i, err))
		}
	}
	if len(parts) > 0 {
		return fmt.Errorf("sequential errors: %s", strings.Join(parts, "; "))
	}
	return nil
}

package util

import (
	"context"
	"escope/internal/config"
	"escope/internal/constants"
	"time"
)

func CreateTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(getDefaultTimeout())*time.Second)
}

func ExecuteWithTimeout[T any](fn func() (T, error)) (T, error) {
	ctx, cancel := CreateTimeoutContext()
	defer cancel()

	resultChan := make(chan struct {
		result T
		err    error
	}, 1)

	go func() {
		result, err := fn()
		resultChan <- struct {
			result T
			err    error
		}{result, err}
	}()

	select {
	case result := <-resultChan:
		return result.result, result.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

func getDefaultTimeout() int {
	timeout, err := config.GetConnectionTimeout()
	if err != nil {
		return constants.DefaultTimeout
	}
	return timeout
}

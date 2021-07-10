package retry

import "time"

type RetryFunc func() error

func Retry(fn func() error, attempts uint) (err error) {
	return RetryDelay(fn, attempts, 0)
}

func RetryDelay(fn func() error, attempts uint, delay time.Duration) (err error) {
	for i := attempts; i > 0; i-- {
		err = fn()
		if err == nil {
			break
		}
		time.Sleep(delay)
	}

	return
}

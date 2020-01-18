package retry

import (
	"time"
)

type RetryFunc func(int) error

// Do retries calling function f n-times.
// It returns an error if none of the tries succeeds.
func Do(n int, f RetryFunc) (err error) {
	for i := 0; i < n; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
	}
	return err
}

// DoSleep retries calling function f n-times and sleeps for d after each call.
// It returns an error if none of the tries succeeds.
func DoSleep(n int, d time.Duration, f RetryFunc) (err error) {
	for i := 0; i < n; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
		time.Sleep(d)
	}
	return err
}

// Forever keeps trying to call function f until it succeeds.
func Forever(f RetryFunc) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
	}
}

// ForeverSleep keeps trying to call function f until it succeeds, and sleeps after each failure.
func ForeverSleep(d time.Duration, f RetryFunc) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
		time.Sleep(d)
	}
}

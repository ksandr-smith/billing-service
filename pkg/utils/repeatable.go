package utils

import "time"

func DoWithTries(f func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = f(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}

		return nil
	}

	return
}

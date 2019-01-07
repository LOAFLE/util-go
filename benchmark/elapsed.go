package benchmark

import (
	"time"
)

func Elapsed() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}

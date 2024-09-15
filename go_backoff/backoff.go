package main

import (
	"math"
	"math/rand/v2"
	"time"
)

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

// ExponentialBackoff - min(cap, 2^attempt * base))
func ExponentialBackoff(attempt int, base, cap time.Duration) time.Duration {
	if attempt <= 0 {
		return base
	}
	if base <= 0 {
		base = 1
	}
	if cap <= 0 {
		cap = math.MaxInt64 // unlimited
	}
	if attempt > 63 || base > (cap>>uint(attempt)) {
		return cap
	}
	return base * time.Duration(1<<uint(attempt))
}

// ExponentialBackoffFullJitter - random_between(0, exp_backoff)
func ExponentialBackoffFullJitter(attempt int, base, cap time.Duration) time.Duration {
	backoff := ExponentialBackoff(attempt, base, cap)
	return time.Duration(rand.Int64N(int64(backoff)))
}

// ExponentialBackoffEqualJitter - exp_backoff / 2 + random_between(0, exp_backoff / 2)
func ExponentialBackoffEqualJitter(attempt int, base, cap time.Duration) time.Duration {
	backoff := ExponentialBackoff(attempt, base, cap)
	if backoff < 2 {
		return backoff
	}
	return backoff/2 + time.Duration(rand.Int64N(int64(backoff/2)))
}

// ExponentialBackoffDecorrelatedJitter - min(cap, random_between(base, exp_backoff * 3))
func ExponentialBackoffDecorrelatedJitter(attempt int, base, cap time.Duration) time.Duration {
	if cap <= 0 {
		cap = math.MaxInt64 // unlimited
	}
	backoff := ExponentialBackoff(attempt-1, base, cap)
	if backoff > math.MaxInt64/3 {
		backoff = math.MaxInt64 / 3
	}
	return min(cap, time.Duration(rand.Int64N(int64(backoff*3-base)))+base)
}

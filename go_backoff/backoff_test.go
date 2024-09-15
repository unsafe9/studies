package main

import (
	"testing"
)

func TestExponentialBackoff(t *testing.T) {
	for i := 0; i < 66; i++ {
		t.Logf("retry:%d\tbackoff:%d", i, ExponentialBackoff(i, 1, 0))
	}
}

func TestExponentialBackoffFullJitter(t *testing.T) {
	for i := 0; i < 66; i++ {
		t.Logf("retry:%d\tbackoff:%d", i, ExponentialBackoffFullJitter(i, 1, 0))
	}
}

func TestExponentialBackoffEqualJitter(t *testing.T) {
	for i := 0; i < 66; i++ {
		t.Logf("retry:%d\tbackoff:%d", i, ExponentialBackoffEqualJitter(i, 1, 0))
	}
}

func TestExponentialBackoffDecorrelatedJitter(t *testing.T) {
	for i := 0; i < 66; i++ {
		t.Logf("retry:%d\tbackoff:%d", i, ExponentialBackoffDecorrelatedJitter(i, 1, 0))
	}
}

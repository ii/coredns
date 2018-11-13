// Package jitter contains various function that return a jittered value.
package jitter

import (
	"math/rand"
	"time"
)

// DurationMillisecond returns a random duration between [0,n) * time.Millisecond
func DurationMillisecond(n int) time.Duration {
	r := rand.Intn(n)
	return time.Duration(r) * time.Millisecond
}

// Add returns t with a random fraction of d added to it.
func Add(t time.Time, d time.Duration) time.Time {
	r := rand.Float64() * float64(d)
	return t.Add(time.Duration(r))
}

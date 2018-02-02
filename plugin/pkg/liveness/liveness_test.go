package liveness

import (
	"testing"
	"time"
)

// TODO(miek): non-sensical test
func TestLiveness(t *testing.T) {
	upfunc := func(s string) bool {
		println("hello", s)
		return true
	}
	pr := New(5 * time.Millisecond)
	defer pr.Stop()
	pr.Start("localhost:8043")
	pr.Do(upfunc)
}

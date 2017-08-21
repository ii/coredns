package kubernetes

import (
	"fmt"
	"testing"
)

func TestAxfr(t *testing.T) {
	k := New([]string{"cluster.local."})
	k.APIConn = &APIConnServeTest{}

	x := NewXfr(k)
	x.services("example.org.")
	fmt.Println()
}

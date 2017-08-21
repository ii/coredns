package kubernetes

import (
	"fmt"
	"testing"
)

func TestAxfr(t *testing.T) {
	k := New([]string{"cluster.local."})
	k.APIConn = &APIConnServeTest{}

	x := &Xfr{k}
	ksrv, _ := x.services("example.org.")
	records := x.getRecordsForK8sItems(ksrv, nil, "example.org.")

	for _, r := range records {
		fmt.Printf("%#v\n", r)
	}
}

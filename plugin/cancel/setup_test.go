package cancel

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", `cancel`)
	if err := setup(c); err != nil {
		t.Fatalf("Test 1, expected no errors, but got: %q", err)
	}

	c = caddy.NewTestController("dns", `cancel aaa`)
	if err := setup(c); err != nil {
		t.Fatalf("Test 2, expected no errors, but got: %q", err)
	}
}

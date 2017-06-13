package debug

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/coredns/coredns/core/dnsserver"

	"github.com/mholt/caddy"
)

func TestDebug(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// Predefined error substrings
	parseErrContent := "Parse error:"

	tests := []struct {
		input     string
		shouldErr bool
	}{
		// positive
		{
			`debug`, false,
		},
		// negative
		{
			`debug off`, true,
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		err := setup(c)
		cfg := dnsserver.GetConfig(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Expected no error but found one for input %s. Error was: %v", i, test.input, err)
			}
		}
	}
}

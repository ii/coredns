package debug

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/mholt/caddy"
)

func TestDebug(t *testing.T) {
	log.SetOutput(ioutil.Discard)

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

package file

import (
	"strings"

	"github.com/coredns/coredns/plugin/pkg/fuzz"
	"github.com/coredns/coredns/plugin/test"
)

func Fuzz(data []byte) int {
	zone, err := Parse(strings.NewReader(dbMiekNL), testzone, "stdin", 0)
	if err != nil {
		t.Fatalf("expect no error when reading zone, got %q", err)
	}
	f := File{Next: test.ErrorHandler(), Zones: Zones{Z: map[string]*Zone{testzone: zone}, Names: []string{testzone}}}

	return fuzz.Do(f, data)
}

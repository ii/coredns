package test

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/coredns/coredns/core/dnsserver"

	"github.com/mholt/caddy"
)

// As we use the filesystem as-is, these files need to exist ON DISK for the readme test to work. This is especially
// useful for the *file* and *dnssec* plugins as their Corefiles are now tested as well. We create empty files in the
// current dir for all these, meaning the example READMEs my use relative path in their READMEs.
var filesToTouch = []string{
	"Kexample.org.+013+45330.key",
	"Kexample.org.+013+45330.private",
	"Kcluster.local+013+45129.key",
	"Kcluster.local+013+45129.private",
}

// TestReadme parses all README.mds of the plugins and checks if every example Corefile
// actually works. Each corefile snippet is only used if the language is set to 'corefile':
//
// ~~~ corefile
// . {
//	# check-this-please
// }
// ~~~
func TestReadme(t *testing.T) {
	port := 30053
	caddy.Quiet = true
	dnsserver.Quiet = true

	touch(filesToTouch)
	defer remove(filesToTouch)

	log.SetOutput(ioutil.Discard)

	middle := filepath.Join("..", "plugin")
	dirs, err := ioutil.ReadDir(middle)
	if err != nil {
		t.Fatalf("Could not read %s: %q", middle, err)
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		readme := filepath.Join(middle, d.Name())
		readme = filepath.Join(readme, "README.md")

		inputs, err := corefileFromReadme(readme)
		if err != nil {
			continue
		}

		// Test each snippet.
		for _, in := range inputs {
			dnsserver.Port = strconv.Itoa(port)
			server, err := caddy.Start(in)
			if err != nil {
				t.Errorf("Failed to start server with %s, for input %q:\n%s", readme, err, in.Body())
			}
			server.Stop()
			port++
		}
	}
}

// corefileFromReadme parses a readme and returns all fragments that
// have ~~~ corefile (or ``` corefile).
func corefileFromReadme(readme string) ([]*Input, error) {
	f, err := os.Open(readme)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	input := []*Input{}
	corefile := false
	temp := ""

	for s.Scan() {
		line := s.Text()
		if line == "~~~ corefile" || line == "``` corefile" {
			corefile = true
			continue
		}

		if corefile && (line == "~~~" || line == "```") {
			// last line
			input = append(input, NewInput(temp))

			temp = ""
			corefile = false
			continue
		}

		if corefile {
			temp += line + "\n" // readd newline stripped by s.Text()
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}
	return input, nil
}

func touch(names []string) {
	for i := range names {
		os.Create(names[i])
	}
}

func remove(names []string) {
	for i := range names {
		os.Remove(names[i])
	}
}

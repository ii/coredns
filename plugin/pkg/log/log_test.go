package log

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDebugLog(t *testing.T) {
	var f bytes.Buffer
	log.SetOutput(&f)

	// D == false
	Debug("debug")
	if x := f.String(); x != "" {
		t.Errorf("Expected no debug logs, got %s", x)
	}

	D = true
	Debug("debug")
	if x := f.String(); !strings.Contains(x, debug+"debug") {
		t.Errorf("Expected debug log to be %s, got %s", debug+"debug", x)
	}
}

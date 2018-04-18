package log

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
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

func TestDebugx(t *testing.T) {
	var f bytes.Buffer
	log.SetOutput(&f)

	D = true

	Debugf("%s", "debug")
	if x := f.String(); !strings.Contains(x, debug+"debug") {
		t.Errorf("Expected debug log to be %s, got %s", debug+"debug", x)
	}

	Debugln("debug")
	if x := f.String(); !strings.Contains(x, debug+"debug\n") {
		t.Errorf("Expected debug log to be %s, got %s", debug+"debug", x)
	}

}

func TestPrint(t *testing.T) {
	var f bytes.Buffer
	log.SetOutput(&f)

	Print("debug")
	if x := f.String(); !strings.Contains(x, "debug") {
		t.Errorf("Expected log to be %s, got %s", "debug", x)
	}
}

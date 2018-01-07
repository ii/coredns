package fall

import "testing"

func TestIsNil(t *testing.T) {
	var f *F
	if !f.IsNil() {
		t.Errorf("F should be nil")
	}
}

func TestIsZero(t *testing.T) {
	f := New()
	if !f.IsZero() {
		t.Errorf("F should be zero")
	}
}

func TestThrough(t *testing.T) {
	if !Example.Through("example.org.") {
		t.Errorf("example.org. should fall through")
	}
	if Example.Through("example.net.") {
		t.Errorf("example.net. should not fall through")
	}

	fall := New()
	if !fall.Through("example.org.") {
		t.Errorf("example.org. should fall through")
	}
}

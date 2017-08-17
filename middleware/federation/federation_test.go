package federation

import "testing"

func TestIsNameFederation(t *testing.T) {
	tests := []struct {
		fed           string
		qname         string
		expectedLabel string
	}{
		{"prod", "nginx.mynamespace.prod.svc.example.com.", "prod"},
		{"prod", "nginx.mynamespace.staging.svc.example.com.", ""},
		{"prod", "nginx.mynamespace.example.com.", ""},
		{"prod", "example.com.", ""},
		{"prod", "com.", ""},
	}

	fed := New()
	for i, tc := range tests {
		fed.f[tc.fed] = "test-name"
		x, y := fed.isNameFederation(tc.qname, "example.com.")
		if x != tc.expectedLabel {
			t.Errorf("Test %d, failed to get label, expected %s, got %s", i, tc.expectedLabel, x)
		}
	}
}

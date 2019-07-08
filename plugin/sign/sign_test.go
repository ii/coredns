package sign

import "testing"

func TestSign(t *testing.T) {
	s := &Sign{nil, 0, 0, "db.miek.nl"}
	s.Sign("miek.nl.")
}

// Package fall handles the fallthrough logic used in plugins that support it.
package fall

import (
	"github.com/coredns/coredns/plugin"
)

// F can be nil to allow for no fallthrough, empty allow all zones to fallthrough or
// contain a zone list that is checked.
type F []string

// New returns a new F.
func New() *F { return new(F) }

func (f *F) Through(qname string) bool {
	if f == nil {
		return false
	}
	if len(*f) == 0 {
		return true
	}
	zone := plugin.Zones(*f).Matches(qname)
	return zone != ""
}

// SetZones will set zones in f.
func (f *F) SetZones(zones []string) {
	for i := range zones {
		zones[i] = plugin.Host(zones[i]).Normalize()
	}
	*f = zones
}

// Example returns an F with example.org. as the zone name.
var Example = func() *F {
	f := F([]string{"example.org."})
	return &f
}()

// Zero returns a zero valued F.
var Zero = func() *F {
	f := F([]string{})
	return &f
}

func (f *F) Len() int { return len(*f) }

func (f *F) IsNil() bool { return f == nil }

func (f *F) IsZero() bool {
	if f == nil {
		return false
	}
	return len(*f) == 0
}

// Equal returns true if f and g are equal. Only useful in tests, The (possible) zones
// are *not* checked.
func (f *F) Equal(g *F) bool {
	if f.IsNil() {
		if g.IsNil() {
			return true
		}
		return false
	}
	if f.IsZero() {
		if g.IsZero() {
			return true
		}
	}
	if f.Len() != g.Len() {
		return false
	}
	return true
}

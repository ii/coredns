package tree

import "github.com/miekg/dns"

// Walk performs fn on all values stored in the tree. A boolean is returned indicating whether the
// Walk was interrupted by an fn returning true. If fn alters stored values' sort
// relationships, future tree operation behaviors are undefined.
func (t *Tree) Walk(fn func(e *Elem, rrs map[uint16][]dns.RR) bool) bool {
	if t.Root == nil {
		return false
	}
	return t.Root.walk(fn)
}

func (n *Node) walk(fn func(e *Elem, rrs map[uint16][]dns.RR) bool) (done bool) {
	if n.Left != nil {
		done = n.Left.walk(fn)
		if done {
			return
		}
	}
	done = fn(n.Elem, n.Elem.m)
	if done {
		return
	}
	if n.Right != nil {
		done = n.Right.walk(fn)
	}
	return
}

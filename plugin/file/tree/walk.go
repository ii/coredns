package tree

// Walk performs fn on all values stored in the tree. A boolean is returned indicating whether the
// Walk was interrupted by an fn returning true. If fn alters stored values' sort
// relationships, future tree operation behaviors are undefined.
func (t *Tree) Walk(fn func(e *Elem) bool) bool {
	if t.Root == nil {
		return false
	}
	return t.Root.walk(fn)
}

func (n *Node) walk(fn func(e *Elem) bool) (done bool) {
	if n.Left != nil {
		done = n.Left.walk(fn)
		if done {
			return
		}
	}
	done = fn(n.Elem)
	if done {
		return
	}
	if n.Right != nil {
		done = n.Right.walk(fn)
	}
	return
}

package metrics

// addrs keeps track on which addrs we listen, so we only start one listener, is
// prometheus is used in multiple Server Blocks.
type addrs struct {
	a map[string]int
}

var uniqAddr addrs

func (a *addrs) setAddress(addr string) {
	// If already there and set to done, we've already started this listener.
	if a.a[addr] == done {
		return
	}
	a.a[addr] = todo
}

// forEachTodo iterates for a and executes f for each element that is 'todo' and sets it to 'done'.
func (a *addrs) forEachTodo(f func() error) {
	for a, v := range uniqAddr.a {
		if v == todo {
			f()
		}
		uniqAddr.a[a] = done
	}
}

// forEachDone iterates for a and executes f for each element that is 'done' and sets it to 'todo'.
func (a *addrs) forEachDone(f func() error) {
	for a, v := range uniqAddr.a {
		if v == done {
			f()
		}
		uniqAddr.a[a] = todo
	}
}

const (
	todo = 1
	done = 2
)

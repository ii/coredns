package forward

func (h host) Check() {
	h.Lock()

	if h.Checking {
		h.Unlock()
		return
	}

	h.Checking = true
	h.Unlock()

	return
}

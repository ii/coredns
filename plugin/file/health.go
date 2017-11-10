package file

// Health implements the health.Healther interface.
func (f File) Health() bool {
	f.Lock()
	defer f.Unlock()
	return f.synced
}

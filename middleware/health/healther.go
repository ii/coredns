package health

type Healther interface {
	// Health returns a boolean indicating the health status of a middleware.
	// False indicates unhealthy.
	Health() bool
}

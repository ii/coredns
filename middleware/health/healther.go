package health

// Healther interface needs to be implemented by each middleware willing to
// provide healthhceck information to the health middleware. As a second step
// the middleware needs to registered against the health middleware, by addding
// it to healthers map. Note this method should return quickly, i.e. just
// checking a boolean status, as it is called every second from the health
// middleware.
type Healther interface {
	// Health returns a boolean indicating the health status of a middleware.
	// False indicates unhealthy.
	Health() bool
}

// Middleware that implements the Healther interface.
var healthers = map[string]bool{
	"erratic": true,
}

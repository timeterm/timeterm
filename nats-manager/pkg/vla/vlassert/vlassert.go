package vlassert

import "gitlab.com/timeterm/timeterm/nats-manager/pkg/vla"

// TestingT is a minimal interface which is implemented by testing.T.
// Is contains the required methods for the assertions to function.
type TestingT interface {
	Errorf(format string, args ...interface{})
	Helper()
}

// IsRegistered checks if a route with method 'method' and path 'path' is registered on the router 'r'.
func IsRegistered(t TestingT, r vla.Router, method, path string) bool {
	t.Helper()

	for _, route := range r.Routes() {
		if route.Path() == path && route.Method() == method {
			return true
		}
	}

	t.Errorf("Route %s %s is not registered on the router", method, path)

	return false
}

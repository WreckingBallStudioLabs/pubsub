//////
// Shared consts, and vars.
//////

package shared

import (
	"net/http"
	"time"
)

//////
// Environments.
//////

const (
	// Development is the development environment.
	Development = "development"

	// Integration is the integration environment.
	Integration = "integration"

	// Production is the production environment.
	Production = "production"

	// Staging is the staging environment.
	Staging = "staging"

	// Testing is the testing environment.
	Testing = "testing"
)

//////
// Timeout.
//////

// Timeout is the default timeout.
const Timeout = 30 * time.Second

// HTTPMethod is the HTTP method.
type HTTPMethod string

const (
	// MethodGet is the HTTP GET method.
	MethodGet HTTPMethod = http.MethodGet

	// MethodPost is the HTTP POST method.
	MethodPost HTTPMethod = http.MethodPost

	// MethodPut is the HTTP PUT method.
	MethodPut HTTPMethod = http.MethodPut

	// MethodPatch is the HTTP PATCH method.
	MethodPatch HTTPMethod = http.MethodPatch

	// MethodDelete is the HTTP DELETE method.
	MethodDelete HTTPMethod = http.MethodDelete
)

func (m HTTPMethod) String() string {
	return string(m)
}

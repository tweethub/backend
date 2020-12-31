package v1

import (
	"fmt"

	"github.com/go-chi/chi"
)

const (
	serviceName = "tweets"
	// version is the API version.
	version = "v1"
)

// setupRoutes setups the REST API routes.
func (srv *Server) setupRoutes() {
	srv.router.Route(fmt.Sprintf("/%s/%s", serviceName, version), func(r chi.Router) {
		r.Get("/user/{user}", srv.ListTweets)
	})
}

package v1

import (
	"fmt"

	"github.com/go-chi/chi"
)

const (
	serviceName = "statistics"
	// version is the API version.
	version = "v1"
	userKey = "user"
)

// setupRoutes setups the REST API routes.
func (srv *Server) setupRoutes() {
	srv.router.Route(fmt.Sprintf("/%s/%s", serviceName, version), func(r chi.Router) {
		r.Route(fmt.Sprintf("/user/{%s}", userKey), func(r chi.Router) {
			r.Route("/relevance", func(r chi.Router) {
				r.Get("/", srv.ListRelevance)
				r.Get("/10d", srv.ListRelevance10d)
				r.Get("/1m", srv.ListRelevance1m)
				r.Get("/3m", srv.ListRelevance3m)
				r.Get("/6m", srv.ListRelevance6m)
				r.Get("/1y", srv.ListRelevance1y)
				r.Get("/5y", srv.ListRelevance5y)
				r.Get("/all-time", srv.ListRelevanceAllTime)
			})
		})
	})
}
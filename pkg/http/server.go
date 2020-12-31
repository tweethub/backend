package http

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opencensus.io/plugin/ochttp"
)

const routerTimeout = 42 * time.Second

// ContentType is a middleware setting the content type.
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// NewRouter returns chi router.
func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(routerTimeout))
	router.Use(ContentType)
	return router
}

func ListenAndServe(handler http.Handler, config Config) chan error {
	errorStream := make(chan error)

	go func() {
		var err error

		addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
		if config.CertFile == "" {
			err = http.ListenAndServe(
				addr,
				&ochttp.Handler{
					Handler: handler,
				},
			)
		} else {
			err = http.ListenAndServeTLS(
				addr,
				config.CertFile,
				config.KeyFile,
				&ochttp.Handler{
					Handler: handler,
				},
			)
		}
		errorStream <- err
		close(errorStream)
	}()
	return errorStream
}

// IsBrokenPipe checks if the error is a broken pipe error.
func IsBrokenPipe(e error) bool {
	return errors.Is(e, syscall.EPIPE)
}

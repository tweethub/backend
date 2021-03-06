package v1

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/tweethub/backend/cmd/statistics/storage"
	"github.com/tweethub/backend/pkg/http"
	"github.com/tweethub/backend/pkg/json"
	"github.com/tweethub/backend/pkg/service"
	"github.com/tweethub/backend/pkg/tweets"
	"go.uber.org/zap"
)

// Server represents a web server.
type Server struct {
	config     http.Config
	router     *chi.Mux
	logger     *zap.Logger
	twtsClient *tweets.Client
	db         storage.Statistical
}

// NewServer returns a server.
func NewServer(db storage.Statistical, twtsClient *tweets.Client, config http.Config, logger *zap.Logger) *Server {
	srv := &Server{
		router:     http.NewRouter(),
		db:         db,
		twtsClient: twtsClient,
		config:     config,
		logger:     logger,
	}
	srv.setupRoutes()
	return srv
}

func (srv *Server) Start(ctx context.Context) {
	srv.logger.Info(service.InitRESTServer)

	errStream := http.ListenAndServe(srv.router, srv.config)

	select {
	case err := <-errStream:
		srv.logger.Error("Listen and serve failed", zap.Error(err))
	default:
		srv.logger.Info(service.RunningRESTServer,
			zap.String("host", srv.config.Host),
			zap.String("port", srv.config.Port))

		<-ctx.Done()
	}
}

func (srv *Server) response(encoder *json.Encoder) {
	if err := encoder.Encode(); err != nil {
		if http.IsBrokenPipe(err) {
			return
		}
		srv.logger.Error("Encoding message response failed", zap.Error(err))
		encoder.SetStatusInternalServerError()
		return
	}
}

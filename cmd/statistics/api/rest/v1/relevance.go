package v1

import (
	"context"
	goerrors "errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	v1 "github.com/tweethub/backend/api/services/statistics/v1"
	"github.com/tweethub/backend/cmd/statistics/generator"
	"github.com/tweethub/backend/cmd/statistics/storage"
	"github.com/tweethub/backend/pkg/json"
	"github.com/tweethub/backend/pkg/time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	allTimeDuration = -1
	defaultSeries   = 20
)

// ListRelevance handles the relevance statistic endpoint.
// @Summary User relevance
// @Description Get relevance statistics of a twitter user based on
// @Description time span, graph series, before and after a specific date.
// @Tags Statistics
// @ID list-relevance
// @Produce json
// @Param user path string true "User"
// @Param before_time query string false "statistics before specific date"
// @Param after_time query string false "statistics after specific date"
// @Param series query int false "distribution of data" mininum(1) default(20)
// @Success 200 {object} v1.RelevanceResponse
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance [get]
func (srv *Server) ListRelevance(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	defer srv.response(encoder)

	if err := r.ParseForm(); err != nil {
		srv.logger.Debug("Parsing form failed", zap.Error(err))
		encoder.SetBody(v1.NewInvalidRequestError()).SetStatusBadRequest()
		return
	}

	rlvURLValues := v1.RelevanceURLValues{}
	if err := schema.NewDecoder().Decode(&rlvURLValues, r.Form); err != nil {
		srv.logger.Debug("Decoding URL values failed", zap.Error(err))
		encoder.SetBody(v1.NewInvalidRequestError()).SetStatusBadRequest()
		return
	}

	if err := ValidateRelevanceURLValues(rlvURLValues); err != nil {
		encoder.SetBody(v1.NewErrorResponse(err.Error())).SetStatusBadRequest()
		return
	}

	user := chi.URLParam(r, userKey)
	ctx := r.Context()

	response, err := srv.generateRelevance(ctx, user, rlvURLValues)
	switch {
	case goerrors.Is(err, mongo.ErrNoDocuments) || response == nil:
		encoder.SetBody(v1.NewNotFoundError()).SetStatusNotFound()
	case err == nil:
		encoder.SetBody(response).SetStatusOK()
	default:
		srv.logger.Error("List relevance failed", zap.Error(err))
		encoder.SetBody(v1.NewInternalServerError()).SetStatusInternalServerError()
	}
}

// ListRelevance10d handles the 10 days relevance statistic endpoint.
// @Summary User relevance for the past 10 days
// @Description Get relevance statistics of a twitter user for the past 10 days.
// @Tags Statistics
// @ID list-relevance-10d
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/10d [get]
func (srv *Server) ListRelevance10d(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, 10*time.Day) // nolint:gomnd
}

// ListRelevance1m handles the 1 month relevance statistic endpoint.
// @Summary User relevance for the past 1 month
// @Description Get relevance statistics of a twitter user for the past 1 month.
// @Tags Statistics
// @ID list-relevance-1m
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/1m [get]
func (srv *Server) ListRelevance1m(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, time.Month)
}

// ListRelevance3m handles the 3 months relevance statistic endpoint.
// @Summary User relevance for the past 3 months
// @Description Get relevance statistics of a twitter user for the past 3 months.
// @Tags Statistics
// @ID list-relevance-3m
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/3m [get]
func (srv *Server) ListRelevance3m(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, 3*time.Month) // nolint:gomnd
}

// ListRelevance6m handles the 6 months relevance statistic endpoint.
// @Summary User relevance for the past 6 months
// @Description Get relevance statistics of a twitter user for the past 6 months.
// @Tags Statistics
// @ID list-relevance-6m
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/6m [get]
func (srv *Server) ListRelevance6m(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, 6*time.Month) // nolint:gomnd
}

// ListRelevance1y handles the 1 year relevance statistic endpoint.
// @Summary User relevance for the past 1 year
// @Description Get relevance statistics of a twitter user for the past 1 year.
// @Tags Statistics
// @ID list-relevance-1y
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/1y [get]
func (srv *Server) ListRelevance1y(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, time.Year)
}

// ListRelevance5y handles the 5 years relevance statistic endpoint.
// @Summary User relevance for the past 5 years
// @Description Get relevance statistics of a twitter user for the past 5 years.
// @Tags Statistics
// @ID list-relevance-5y
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/5y [get]
func (srv *Server) ListRelevance5y(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, 5*time.Year) // nolint:gomnd
}

// ListRelevanceAllTime handles the all time relevance statistic endpoint.
// @Summary All time user relevance
// @Description Get all time relevance statistics of a twitter user.
// @Tags Statistics
// @ID list-relevance-all-time
// @Produce json
// @Param user path string true "User"
// @Success 200 {object} v1.RelevanceResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /statistics/v1/user/{user}/relevance/all-time [get]
func (srv *Server) ListRelevanceAllTime(w http.ResponseWriter, r *http.Request) {
	srv.listRelevance(w, r, allTimeDuration)
}

func (srv *Server) listRelevance(w http.ResponseWriter, r *http.Request, timeSpan time.Duration) {
	encoder := json.NewEncoder(w)
	defer srv.response(encoder)

	user := chi.URLParam(r, userKey)
	ctx := r.Context()

	opts := storage.RelevanceOptions{}
	if timeSpan == allTimeDuration {
		opts.TimeSpan = nil
	} else {
		opts.TimeSpan = &timeSpan
	}

	response, err := srv.db.Relevance(ctx, user, opts)
	switch {
	case goerrors.Is(err, mongo.ErrNoDocuments) || response == nil:
		encoder.SetBody(v1.NewNotFoundError()).SetStatusNotFound()
	case err == nil:
		encoder.SetBody(response).SetStatusOK()
	default:
		srv.logger.Error("Getting relevance from database failed",
			zap.Object("options", opts),
			zap.Error(err),
		)
		encoder.SetBody(v1.NewInternalServerError()).SetStatusInternalServerError()
	}
}

// generateRelevance generates relevance.
func (srv *Server) generateRelevance(
	ctx context.Context,
	user string,
	rlvURLValues v1.RelevanceURLValues,
) (*v1.RelevanceResponse, error) {
	// Get the tweets in the time span.
	twts, err := srv.twtsClient.GetTweets(ctx, user, GenerateTweetsURLValues(rlvURLValues))
	if err != nil {
		return nil, errors.Wrap(err, "getting tweets")
	}

	if rlvURLValues.Series == 0 {
		rlvURLValues.Series = defaultSeries
	}

	// Calculate the relevance from the tweets.
	response, err := generator.CalculateRelevance(twts, GenerateTimeFrame(rlvURLValues), rlvURLValues.Series)
	if err != nil {
		return nil, errors.Wrap(err, "calculating relevance")
	}
	return response, nil
}

package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/pkg/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// ListTweets handles the tweets endpoint.
// @Summary User tweets
// @Description Get all the tweets of a twitter user or the tweets in a specific time span.
// @Tags Tweets
// @ID list-tweets
// @Produce json
// @Param user path string true "User"
// @Param before_time query string false "tweets before specific date"
// @Param after_time query string false "tweets after specific date"
// @Success 200 {object} v1.TweetsResponse
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /tweets/v1/user/{user} [get]
func (srv *Server) ListTweets(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	defer srv.response(encoder)

	if err := r.ParseForm(); err != nil {
		srv.logger.Debug("Parsing form failed", zap.Error(err))
		encoder.SetBody(v1.NewInvalidRequestError()).SetStatusBadRequest()
		return
	}

	twtsURLValues := v1.TweetsURLValues{}
	if err := schema.NewDecoder().Decode(&twtsURLValues, r.Form); err != nil {
		srv.logger.Debug("Decoding URL values failed", zap.Error(err))
		encoder.SetBody(v1.NewInvalidRequestError()).SetStatusBadRequest()
		return
	}

	if err := ValidateTweetsURLValues(twtsURLValues); err != nil {
		encoder.SetBody(v1.NewErrorResponse(err.Error())).SetStatusBadRequest()
		return
	}

	ctx := r.Context()
	user := chi.URLParam(r, "user")

	response, err := srv.db.Tweets(ctx, user, GenerateTweetsOptions(twtsURLValues))
	switch {
	case errors.Is(err, mongo.ErrNoDocuments) || response == nil:
		encoder.SetBody(v1.NewNotFoundError()).SetStatusNotFound()
	case err == nil:
		encoder.SetBody(response).SetStatusOK()
	default:
		srv.logger.Error("Getting tweets from database failed", zap.Error(err))
		encoder.SetBody(v1.NewInternalServerError()).SetStatusInternalServerError()
	}
}

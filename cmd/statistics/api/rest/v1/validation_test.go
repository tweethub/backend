package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "github.com/tweethub/backend/api/services/statistics/v1"
	"github.com/tweethub/backend/pkg/twitter"
)

func TestValidateRelevanceURLValues(t *testing.T) {
	type fields struct {
		AfterTime  time.Time
		BeforeTime time.Time
		Series     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "Start date after end date",
			fields: fields{
				Series:     20,
				AfterTime:  twitter.UserFirstTweetDate().Add(10),
				BeforeTime: twitter.UserFirstTweetDate().Add(8),
			},
			wantErr: errInvalidBeforeTimeOrAfterTime,
		},
		{
			name: "Start date in the future",
			fields: fields{
				Series:     20,
				AfterTime:  time.Now().Add(1 * time.Hour),
				BeforeTime: time.Now().Add(2 * time.Hour),
			},
			wantErr: errInvalidAfterTime,
		},
		{
			name: "End date before first tweet",
			fields: fields{
				Series:     20,
				AfterTime:  twitter.UserFirstTweetDate().Add(-2 * time.Hour),
				BeforeTime: twitter.UserFirstTweetDate().Add(-1 * time.Hour),
			},
			wantErr: errInvalidBeforeTime,
		},
		{
			name: "Valid time span and series",
			fields: fields{
				Series:     10,
				AfterTime:  twitter.UserFirstTweetDate(),
				BeforeTime: twitter.UserFirstTweetDate().Add(2 * time.Second),
			},
			wantErr: nil,
		},
		{
			name: "Invalid series",
			fields: fields{
				Series:     -3,
				AfterTime:  twitter.UserFirstTweetDate(),
				BeforeTime: twitter.UserFirstTweetDate().Add(2 * time.Second),
			},
			wantErr: errInvalidSeries,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vls := v1.RelevanceURLValues{
				AfterTime:  tt.fields.AfterTime,
				BeforeTime: tt.fields.BeforeTime,
				Series:     tt.fields.Series,
			}
			err := ValidateRelevanceURLValues(vls)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

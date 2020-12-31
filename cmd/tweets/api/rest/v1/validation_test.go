package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/pkg/twitter"
)

func TestValidateTweetsURLValues(t *testing.T) {
	type fields struct {
		AfterTime  time.Time
		BeforeTime time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "Start date after end date",
			fields: fields{
				AfterTime:  time.Date(0, 0, 0, 0, 2, 0, 0, time.UTC),
				BeforeTime: time.Date(0, 0, 0, 0, 1, 0, 0, time.UTC),
			},
			wantErr: errInvalidBeforeTimeOrAfterTime,
		},
		{
			name: "No error",
			fields: fields{
				AfterTime:  twitter.UserFirstTweetDate(),
				BeforeTime: twitter.UserFirstTweetDate().Add(2 * time.Second),
			},
			wantErr: nil,
		},
		{
			name: "Start date in the future",
			fields: fields{
				AfterTime:  time.Now().Add(1 * time.Hour),
				BeforeTime: time.Now().Add(2 * time.Hour),
			},
			wantErr: errInvalidAfterTime,
		},
		{
			name: "End date before first tweet",
			fields: fields{
				AfterTime:  twitter.UserFirstTweetDate().Add(-2 * time.Hour),
				BeforeTime: twitter.UserFirstTweetDate().Add(-1 * time.Hour),
			},
			wantErr: errInvalidBeforeTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vls := v1.TweetsURLValues{
				AfterTime:  tt.fields.AfterTime,
				BeforeTime: tt.fields.BeforeTime,
			}
			err := ValidateTweetsURLValues(vls)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

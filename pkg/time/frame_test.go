package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculatePeriods(t *testing.T) {
	type args struct {
		fromDate time.Time
		duration time.Duration
		n        int
	}
	tests := []struct {
		name     string
		args     args
		expected []*Frame
	}{
		{
			name: "Test calculating time periods",
			args: args{
				fromDate: time.Date(0, 0, 0, 0, 0, 1, 0, time.UTC),
				duration: 33 * time.Second,
				n:        3,
			},
			expected: []*Frame{
				{
					StartTime: time.Date(0, 0, 0, 0, 0, 1, 0, time.UTC),
					EndTime:   time.Date(0, 0, 0, 0, 0, 34, 0, time.UTC),
				},
				{
					StartTime: time.Date(0, 0, 0, 0, 0, 34, 0, time.UTC),
					EndTime:   time.Date(0, 0, 0, 0, 1, 7, 0, time.UTC),
				},
				{
					StartTime: time.Date(0, 0, 0, 0, 1, 7, 0, time.UTC),
					EndTime:   time.Date(0, 0, 0, 0, 1, 40, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeFrames := GenerateFrames(tt.args.fromDate, tt.args.duration, tt.args.n)

			require.Equal(t, tt.expected, timeFrames)
		})
	}
}

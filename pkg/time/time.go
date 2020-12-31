package time

import (
	"sort"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	Day   = 24 * time.Hour
	Month = 30 * Day
	Year  = 12 * Month
)

// NowUTC returns the time now in UTC.
func NowUTC() time.Time {
	return time.Now().UTC()
}

type Duration = time.Duration

type Durations []time.Duration

// Sort sorts the durations.
func (d Durations) Sort() {
	sort.Slice(d, func(i, j int) bool {
		return d[i] < d[j]
	})
}

func (d Durations) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, duration := range d {
		enc.AppendDuration(duration)
	}
	return nil
}
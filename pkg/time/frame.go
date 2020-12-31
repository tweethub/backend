package time

import (
	"fmt"
	"time"
)

// Frame represents a period of time.
type Frame struct {
	StartTime time.Time
	EndTime   time.Time
}

// ToUnixString converts time frame to unix string.
func (f Frame) ToUnixString() string {
	return fmt.Sprintf("%d-%d", f.StartTime.Unix(), f.EndTime.Unix())
}

// ToDurationString converts time frame to duration string.
func (f Frame) ToDurationString() string {
	return f.Duration().String()
}

// Duration returns the time frame duration.
func (f Frame) Duration() Duration {
	return f.EndTime.Sub(f.StartTime)
}

// UTC sets the time frame to UTC.
func (f *Frame) UTC() {
	f.StartTime = f.StartTime.UTC()
	f.EndTime = f.EndTime.UTC()
}

// GenerateFrames generates time frames.
func GenerateFrames(start time.Time, duration time.Duration, n int) []*Frame {
	frames := make([]*Frame, n)

	for i := 0; i < n; i++ {
		endTime := start.Add(duration)
		frames[i] = &Frame{
			StartTime: start,
			EndTime:   endTime,
		}
		start = endTime
	}
	return frames
}

package time

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortDurations(t *testing.T) {
	tests := []struct {
		name      string
		durations Durations
		expected  Durations
	}{
		{
			name:      "Test sorting durations",
			durations: Durations{6, 1, 2, 3, 5, 4},
			expected:  Durations{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.durations.Sort()

			require.Equal(t, tt.expected, tt.durations)
		})
	}
}

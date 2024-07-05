package haar

import (
	"fmt"
	"testing"

	"github.com/jnafolayan/sip/pkg/signal"
)

func TestHaarTransform(t *testing.T) {
	tests := []struct {
		level    int
		signal   signal.Signal2D
		expected signal.Signal2D
	}{
		{
			1,
			signal.Signal2D{[]float32{3, 5, 4, 8}},
			signal.Signal2D{[]float32{4, 6, -1, -2}},
		},
	}

	hw := &HaarWavelet{Level: 1}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level %d decomposition", tt.level), func(t *testing.T) {
			ts := hw.Transform(tt.signal)
			if !ts.Equal(tt.expected) {
				tsStr := ts.String(ts.Bounds())
				expStr := tt.expected.String(tt.expected.Bounds())
				t.Errorf("Transformed signal was incorrect, got:\n%s\nwant:\n%s", tsStr, expStr)
			}
		})
	}
}

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
			signal.Signal2D{{3, 5, 4, 8}},
			signal.Signal2D{{4, 6, -1, -2}, {0, 0, 0, 0}},
		},
		{
			2,
			signal.Signal2D{{3, 5, 4, 8}},
			signal.Signal2D{{5, -1, -1, -2}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level %d decomposition", tt.level), func(t *testing.T) {
			hw := &HaarWavelet{Level: tt.level}
			ts := hw.Transform(tt.signal)
			if !ts.Equal(tt.expected) {
				tsStr := ts.String(ts.Bounds())
				expStr := tt.expected.String(tt.expected.Bounds())
				t.Errorf("Transformed signal was incorrect, got:\n%s\nwant:\n%s", tsStr, expStr)
			}
		})
	}
}

func TestInverseHaarTransform(t *testing.T) {
	tests := []struct {
		level  int
		signal signal.Signal2D
	}{
		{
			1,
			signal.Signal2D{{3, 5, 4, 8}},
		},
		{
			2,
			signal.Signal2D{{3, 5, 4, 8}},
		},
		{
			1,
			signal.Signal2D{
				{64, 2, 3, 61, 60, 6, 7, 57},
				{9, 55, 54, 12, 13, 51, 50, 16},
				{17, 47, 46, 20, 21, 43, 42, 24},
				{40, 26, 27, 37, 36, 30, 31, 33},
				{32, 34, 35, 29, 28, 38, 39, 25},
				{41, 23, 22, 44, 45, 19, 18, 48},
				{49, 15, 14, 52, 53, 11, 10, 56},
				{8, 58, 59, 5, 4, 62, 63, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level %d decomposition", tt.level), func(t *testing.T) {
			hw := &HaarWavelet{Level: tt.level}
			ts := hw.Transform(tt.signal)
			recon := hw.InverseTransform(ts).Slice(0, 0, len(tt.signal[0]), len(tt.signal))
			if !recon.Equal(tt.signal) {
				reconStr := recon.String(recon.Bounds())
				expStr := tt.signal.String(tt.signal.Bounds())
				t.Errorf("Reconstructed signal was incorrect, got:\n%s\nwant:\n%s", reconStr, expStr)
			}
		})
	}
}

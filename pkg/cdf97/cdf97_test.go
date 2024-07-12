package cdf97

import (
	"fmt"
	"testing"

	"github.com/jnafolayan/sip/pkg/signal"
)

func TestCDF97Transform(t *testing.T) {
	tests := []struct {
		level    int
		signal   signal.Signal2D
		expected signal.Signal2D
	}{
		{
			1,
			signal.Signal2D{{3, 5, 4, 8}},
			signal.Signal2D{{4, 1}, {0, 0}},
		},
		{
			2,
			signal.Signal2D{{3, 5, 4, 8}},
			signal.Signal2D{{4, 1}, {0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level %d decomposition", tt.level), func(t *testing.T) {
			w := &CDF97Wavelet{Level: tt.level}
			ts := w.Transform(tt.signal)
			if !ts.Equal(tt.expected) {
				tsStr := ts.String(ts.Bounds())
				expStr := tt.expected.String(tt.expected.Bounds())
				t.Errorf("Transformed signal was incorrect, got:\n%s\nwant:\n%s", tsStr, expStr)
			}
		})
	}
}

func TestInverseCDF97Transform(t *testing.T) {
	tests := []struct {
		level  int
		signal signal.Signal2D
	}{
		{
			1,
			signal.Signal2D{{3, 5, 4, 8}},
		},
		// {
		// 	2,
		// 	signal.Signal2D{{3, 5, 4, 8}},
		// },
		// {
		// 	1,
		// 	signal.Signal2D{
		// 		{64, 2, 3, 61, 60, 6, 7, 57},
		// 		{9, 55, 54, 12, 13, 51, 50, 16},
		// 		{17, 47, 46, 20, 21, 43, 42, 24},
		// 		{40, 26, 27, 37, 36, 30, 31, 33},
		// 		{32, 34, 35, 29, 28, 38, 39, 25},
		// 		{41, 23, 22, 44, 45, 19, 18, 48},
		// 		{49, 15, 14, 52, 53, 11, 10, 56},
		// 		{8, 58, 59, 5, 4, 62, 63, 1},
		// 	},
		// },
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d - level %d decomposition", i, tt.level), func(t *testing.T) {
			w := &CDF97Wavelet{Level: tt.level}
			ts := w.Transform(tt.signal)
			recon := w.InverseTransform(ts)
			// recon = recon.Slice(0, 0, tt.signal.Bounds().Max.X, tt.signal.Bounds().Max.Y)
			if !recon.Equal(tt.signal) {
				reconStr := recon.String(recon.Bounds())
				expStr := tt.signal.String(tt.signal.Bounds())
				t.Errorf("Reconstructed signal was incorrect, got:\n%s\nwant:\n%s", reconStr, expStr)
			}
		})
	}
}

func TestInverseCDF97Transform1(t *testing.T) {
	s := signal.Signal2D{
		{64, 2, 3, 61, 60, 6, 7, 57},
		{9, 55, 54, 12, 13, 51, 50, 16},
		{17, 47, 46, 20, 21, 43, 42, 24},
		{40, 26, 27, 37, 36, 30, 31, 33},
		{32, 34, 35, 29, 28, 38, 39, 25},
		{41, 23, 22, 44, 45, 19, 18, 48},
		{49, 15, 14, 52, 53, 11, 10, 56},
		{8, 58, 59, 5, 4, 62, 63, 1},
	}

	w := &CDF97Wavelet{Level: 1}
	ts := w.Transform(s)
	recon := w.InverseTransform(ts)

	t.Errorf("\n%s", ts.String(ts.Bounds()))
	t.Errorf("\n%s", recon.String(recon.Bounds()))
}

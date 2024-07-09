package signal

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"slices"
	"strings"
)

type Signal2D [][]float32

func New(width, height int) Signal2D {
	s := make(Signal2D, height)
	for i := 0; i < height; i++ {
		s[i] = make([]float32, width)
	}
	return s
}

func (s Signal2D) Size() (int, int) {
	if len(s) == 0 {
		return 0, 0
	}
	return len(s[0]), len(s)
}

// Clone returns a deep clone of the signal
func (s Signal2D) Clone() Signal2D {
	result := make(Signal2D, len(s))
	for i := range s {
		result[i] = slices.Clone(s[i])
	}
	return result
}

// String produces a string representation of the signal
func (s Signal2D) String(bounds image.Rectangle) string {
	var out strings.Builder

	startY := int(math.Max(0, float64(bounds.Min.Y)))
	endY := int(math.Min(float64(len(s)), float64(bounds.Max.Y)))
	startX := int(math.Max(0, float64(bounds.Min.X)))
	endX := int(math.Min(float64(len(s[0])), float64(bounds.Max.X)))

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if x > startX {
				out.WriteByte('|')
			}
			out.WriteString(fmt.Sprintf(" %5.1f ", s[y][x]))
		}
		out.WriteByte('\n')
	}

	// Remove any trailing newline
	return strings.TrimRight(out.String(), "\n")
}

func (s Signal2D) Image() image.Image {
	w, h := s.Size()
	img := image.NewRGBA(s.Bounds())

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			shade := uint8(s[i][j])
			c := color.RGBA{shade, shade, shade, 255}
			img.Set(j, i, c)
		}
	}

	return img
}

func (s Signal2D) Equal(s2 Signal2D) bool {
	w1, h1 := s.Size()
	w2, h2 := s2.Size()
	if w1 != w2 || h1 != h2 {
		return false
	}

	for i := 0; i < h1; i++ {
		for j := 0; j < w1; j++ {
			if s[i][j] != s2[i][j] {
				return false
			}
		}
	}

	return true
}

// Bounds returns a rectangle defining the size of the signal
func (s Signal2D) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(s[0]), len(s))
}

// PadStyle denotes how new cells in an extended signal will be decided.
type PadStyle int

const (
	// PadZero extends the signal by setting the new cells to zero.
	PadZero PadStyle = iota + 1
	// PadSymmetric extends the signal by repeating the edge cells.
	PadSymmetric
)

// Pad extends the signal to fit a new size fills empty cells using the pad style
// specified in `padStyle`.
func (s Signal2D) Pad(width, height int, padStyle PadStyle) Signal2D {
	oldWidth, oldHeight := s.Size()

	result := New(width, height)

	// Copy existing signal
	var v float32
	for i := range s {
		copy(result[i], s[i])
		// Fill new columns using the pad style specified
		for j := oldWidth; j < width; j++ {
			v = 0
			if padStyle == PadSymmetric {
				v = s[i][oldWidth-1]
			}
			result[i][j] = v
		}
	}

	// Fill new rows using the pad style specified
	for i := oldHeight; i < height; i++ {
		if padStyle == PadSymmetric {
			copy(result[i], result[oldHeight-1])
		}
	}

	return result
}

func (s Signal2D) HardThreshold(offsetX, offsetY, threshold int) Signal2D {
	width, height := s.Size()
	result := s.Clone()
	thresh := float64(threshold)

	for i := offsetY; i < height; i++ {
		for j := offsetX; j < width; j++ {
			if math.Abs(float64(s[i][j])) < thresh {
				result[i][j] = 0
			}
		}
	}

	return result
}

func (s Signal2D) SoftThreshold(offsetX, offsetY, threshold int) Signal2D {
	width, height := s.Size()
	result := s.Clone()
	thresh := float64(threshold)

	var abs float64

	for i := offsetY; i < height; i++ {
		for j := offsetX; j < width; j++ {
			abs = math.Abs(float64(s[i][j]))
			if abs < thresh {
				result[i][j] = 0
			} else {
				result[i][j] = float32(math.Copysign(abs-thresh, float64(s[i][j])))
			}
		}
	}

	return result
}

func (s Signal2D) Slice(x1, y1, x2, y2 int) Signal2D {
	result := New(x2-x1, y2-y1)
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			result[y-y1][x-x1] = s[y][x]
		}
	}
	return result

}
